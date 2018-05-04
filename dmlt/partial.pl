#! /usr/bin/perl

use strict;
use warnings;

my $partFile = "partalExcel.txt";
my @commFiles = (); # list of file name
my %operFiles = (); # hash of list ref

my $dovFile = "策划转表_公共.bat";
my @doconvs = ();

main();
##-- 主过程 --##
sub main
{
	readFileList();
	readDovList();

	if (scalar keys(%operFiles) == 0) {
		convertFiles(@commFiles, "comm");
		print "convert with no oper done! press to finish...\n";
		<STDIN>;
	} else {
		foreach my $oper (keys(%operFiles)) {
			my @listFile = ();
			push(@listFile, @commFiles);
			push(@listFile, @{$operFiles{$oper}});
			convertFiles(@listFile, $oper);

			print "convert $oper done! press to continue...\n";
			my $input = <STDIN>;
			last if $input =~ /^[nq]/i;
		}
	}
}

##-- 子过程 --##

# 解析文件列表
sub readFileList
{
	open(my $fh, "<", $partFile) or die "cannot open $partFile: $!";
	while (<$fh>) {
		chomp;
		next if /^\s*$/;
		next if /^\s*#/;

		# 移除 xls/ 路径之前部分
		s#^.*/xsl/#xsl/#;
		# 将 / 替换为 \
		s#/#\\#g;

		my $oper = "";
		if (/_(.+)\//) {
			$oper = $1;
			if (not exists($operFiles{$oper})) {
				$operFiles{$oper} = [];
			}
			push(@{$operFiles{$oper}}, $_);
		} else {
			push(@commFiles, $_);
		}
	}
}

# 读入批处理中每条转表语句
sub readDovList
{
	open(my $fh, "<", $dovFile) or die "cannot open $dovFile: $!";
	while (<$fh>) {
		chomp;
		next if /^\s*$/;
		next if /^\s*#/;
		next if /^\s*::/;
		next if /^\s*rem/i;
		push(@doconvs, $_);
	}
}


# 转换一个组文件列表
# 参数 $1 Excel文件名列表，$2 标签名
sub convertFiles
{
	my @files = shift;
	my $label = shift;
	$label = "comm" unless defined $label;

	my ($doscmd, $ecode);

	cleanFolder();

	# 复制文件至临时目录 xls_tmp
	foreach my $file (@files) {
		my $dstfile = $file;
		$dstfile =~ s/^.*\\//;
		$doscmd = "xcopy /Y /F $file xls_tmp\\$dstfile";
		$ecode = system($doscmd);
		die "fails to copy $file, exit $ecode" unless $ecode == 0;
		print "[done] $doscmd\n";
	}

	# excel 转 csv
	$doscmd = "x2c\\xls2csv.exe .\\xls_tmp .\\csv x2c.x2c";
	$ecode = system($doscmd);
	die "fails to convert to csv, exit $ecode" unless $ecode == 0;
	print "[done] $doscmd\n";

	# 读取结果 csv 文件列表
	opendir(my $dh, "csv") or die "cannot open dir csv: $!";
	my @csv = readdir($dh);
	my %hascsv = ();
	foreach my $v (@csv) {
		print "$v\n";
		$v =~ s/\.csv//;
		$hascsv{$v} = 1;
	}
	
	# 将 csv 转为 bin
	foreach my $do (@doconvs) {
		my @words = split($do, '\s+');
		foreach my $w (@words) {
			if (exists($hascsv{$w})) {
				$ecode = system($do);
				die "fails to convert to bin, exit $ecode" unless $ecode == 0;
				print "[done] $do\n";
				last;
			}
		}
	}

	# 将结果 bin 文件转存
	mkdir("pack\\$label") or die "cannot create dir $label: $!";
	$doscmd = "xcopy bin\\*.bin pack\\$label";
	$ecode = system($doscmd);
	die "fails to copy out bin, exit $ecode" unless $ecode == 0;
	print "[done] $doscmd\n";
}

# 清理工作目录
sub cleanFolder
{
	my ($doscmd, $ecode);
}




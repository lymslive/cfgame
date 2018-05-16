#! /usr/bin/perl
# 由指定的 excel 文件路径列表，实现部分转表，主要用于运营需求
# 该脚本应该位于 data 数据主目录下的 runenv\ 子目录
# 其余涉及读取的文件都是相对于 data 主目录的路径

use strict;
use warnings;
use Archive::Tar;

my $DEBUG = 1;
my $partFile = "xls\\partialExcel.txt";
my @commFiles = (); # list of file name
my %operFiles = (); # hash of list ref

my $dovFile = "策划转表_公共.bat";
my @doconvs = ();

# 数据主目录在脚本父目录
chdir('..');
die "you may not in the correct data dir" unless (-d 'xls');

main();
##-- 主过程 --##
sub main
{
	readFileList();
	readDovList();

	# seeFileList(); exit 0;

	if (scalar keys(%operFiles) == 0) {
		convertFiles(\@commFiles, "comm");
		elog('prompt',  "convert with no oper done! press to finish...");
	} else {
		foreach my $oper (keys(%operFiles)) {
			my @listFile = ();
			push(@listFile, @commFiles);
			push(@listFile, @{$operFiles{$oper}});
			convertFiles(\@listFile, $oper);

			elog('promt', "convert $oper done! press to continue...");
			my $input = <STDIN>;
			last if $input =~ /^[nq]/i;
		}
	}
	elog('prompt', 'all conversion have done, press to exit');
	<STDIN>;
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
		s#^.*/xls/#xls/#;
		# 将 / 替换为 \
		s#/#\\#g;

		elog("DBG", "read in: $_") if $DEBUG;

		my $oper = "";
		if (/_(.+)\\/) {
			$oper = $1;
			if (not exists($operFiles{$oper})) {
				$operFiles{$oper} = [];
			}
			push(@{$operFiles{$oper}}, $_);
		} else {
			push(@commFiles, $_);
		}
	}
	close($fh);
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
	close($fh);
}


# 转换一个组文件列表
# 参数 $1 Excel文件名列表，$2 标签名
sub convertFiles(\@$)
{
	my $files = shift;
	my $label = shift;
	$label = "comm" unless defined $label;

	elog("INFO", "===== convert oper: $label =====");

	cleanFolder();

	# 复制文件至临时目录 xls_tmp
	foreach my $file (@$files) {
		my $dstfile = $file;
		$dstfile =~ s/^.*\\//;
		dosCmd("xcopy /Y $file xls_tmp");
	}

	# excel 转 csv
	dosCmd("x2c\\xls2csv.exe .\\xls_tmp .\\csv x2c.x2c");

	# 读取结果 csv 文件列表
	opendir(my $dh, "csv") or die "cannot open dir csv: $!";
	my @csv = readdir($dh);
	my %hascsv = ();
	foreach my $v (@csv) {
		next if $v =~ /^\./;
		$v =~ s/\.csv//;
		$hascsv{lc($v)} = 1;
		elog("DBG", "hascsv{$v}\n") if $DEBUG;
	}
	
	# 将 csv 转为 bin
	foreach my $do (@doconvs) {
		my @words = split('\s+', $do);
		foreach my $w (@words) {
			# print "[DBG] word{$w}\n";
			if (exists($hascsv{lc($w)})) {
				dosCmd($do);
				last;
			}
		}
	}

	# 将结果 bin 文件转存并打包
	if (-d "runenv\\$label") {
		dosCmd("rd /S /Q runenv\\$label");
	}
	dosCmd("md runenv\\$label");
	dosCmd("xcopy runenv\\zone_svr\\cfg\\res\\*.bin runenv\\$label");
	my @binfiles = glob("runenv\\zone_svr\\cfg\\res\\*.bin");
	my $tar = "runenv\\operdata_$label.tgz";
	Archive::Tar->create_archive($tar, COMPRESS_GZIP, @binfiles);
	elog("", "result bin saved in $tar");
}

# 清理工作目录
sub cleanFolder
{
	dosCmd('rd /S /Q .\csv') if (-d '.\csv');
	dosCmd('rd /S /Q .\xls_tmp') if (-d '.\xls_tmp');
	dosCmd('rd /S /Q .\bin') if (-d '.\bin');
	dosCmd('rd /S /Q .\runenv\zone_svr\cfg\res') if (-d '.\runenv\zone_svr\cfg\res');

	dosCmd('md .\csv');
	dosCmd('md .\xls_tmp');
	dosCmd('md .\bin');
	dosCmd('md .\runenv\zone_svr\cfg\res');
}

# 调用外部命令
sub dosCmd
{
	my $doscmd = shift;
	my $ecode = system($doscmd);
	die "[fail] exit $ecode; $doscmd" unless $ecode == 0;
	elog('done', $doscmd);
}

# 简单的日志函数，要求前缀与实际消息两个参数，自动添加换行
sub elog($$)
{
	my ($prifex, $msg) = @_;
	$prifex = "INFO" unless length($prifex) > 0;
	chomp($msg);
	print "[$prifex] $msg\n";
}

# 调试查看
sub seeFileList
{
	print "comms files:\n";
	foreach my $file (@commFiles){
		print "  $file\n";
	}

	foreach my $oper (keys(%operFiles)) {
		print "oper $oper files:\n";
		foreach my $file (@{$operFiles{$oper}}){
			print "  $file\n";
		}
	}
}

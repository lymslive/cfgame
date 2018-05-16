#! /usr/bin/perl
# ��ָ���� excel �ļ�·���б�ʵ�ֲ���ת����Ҫ������Ӫ����
# �ýű�Ӧ��λ�� data ������Ŀ¼�µ� runenv\ ��Ŀ¼
# �����漰��ȡ���ļ���������� data ��Ŀ¼��·��

use strict;
use warnings;
use Archive::Tar;

my $DEBUG = 1;
my $partFile = "xls\\partialExcel.txt";
my @commFiles = (); # list of file name
my %operFiles = (); # hash of list ref

my $dovFile = "�߻�ת��_����.bat";
my @doconvs = ();

# ������Ŀ¼�ڽű���Ŀ¼
chdir('..');
die "you may not in the correct data dir" unless (-d 'xls');

main();
##-- ������ --##
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

##-- �ӹ��� --##

# �����ļ��б�
sub readFileList
{
	open(my $fh, "<", $partFile) or die "cannot open $partFile: $!";
	while (<$fh>) {
		chomp;
		next if /^\s*$/;
		next if /^\s*#/;

		# �Ƴ� xls/ ·��֮ǰ����
		s#^.*/xls/#xls/#;
		# �� / �滻Ϊ \
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

# ������������ÿ��ת�����
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


# ת��һ�����ļ��б�
# ���� $1 Excel�ļ����б�$2 ��ǩ��
sub convertFiles(\@$)
{
	my $files = shift;
	my $label = shift;
	$label = "comm" unless defined $label;

	elog("INFO", "===== convert oper: $label =====");

	cleanFolder();

	# �����ļ�����ʱĿ¼ xls_tmp
	foreach my $file (@$files) {
		my $dstfile = $file;
		$dstfile =~ s/^.*\\//;
		dosCmd("xcopy /Y $file xls_tmp");
	}

	# excel ת csv
	dosCmd("x2c\\xls2csv.exe .\\xls_tmp .\\csv x2c.x2c");

	# ��ȡ��� csv �ļ��б�
	opendir(my $dh, "csv") or die "cannot open dir csv: $!";
	my @csv = readdir($dh);
	my %hascsv = ();
	foreach my $v (@csv) {
		next if $v =~ /^\./;
		$v =~ s/\.csv//;
		$hascsv{lc($v)} = 1;
		elog("DBG", "hascsv{$v}\n") if $DEBUG;
	}
	
	# �� csv תΪ bin
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

	# ����� bin �ļ�ת�沢���
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

# ������Ŀ¼
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

# �����ⲿ����
sub dosCmd
{
	my $doscmd = shift;
	my $ecode = system($doscmd);
	die "[fail] exit $ecode; $doscmd" unless $ecode == 0;
	elog('done', $doscmd);
}

# �򵥵���־������Ҫ��ǰ׺��ʵ����Ϣ�����������Զ���ӻ���
sub elog($$)
{
	my ($prifex, $msg) = @_;
	$prifex = "INFO" unless length($prifex) > 0;
	chomp($msg);
	print "[$prifex] $msg\n";
}

# ���Բ鿴
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

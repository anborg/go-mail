::cp testfiles/cayinput566000.csv input/AP566000
::go-mail.exe

:: cp \\markham.ca\apps\eft\test_files\cayinput566000.csv \\markham.ca\apps\eft\input\AP566000

:: Run with uid/pwd
::schtasks /run /tn "\mkm\eft" /s wdnvmhapp98 /u s_application /p app##svc@101
schtasks /run /tn "\mkm\eft" /s wdnvmhapp98 /u p_pmn

:: Run without uid
::schtasks /run /tn "\mkm\eft" /s wdnvmhapp98
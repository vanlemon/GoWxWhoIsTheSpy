sftp root@something << EOF
rename /home/lighthouse/go/src/lmf.mortal.com/GoWxWhoIsTheSpy "/home/lighthouse/go/src/lmf.mortal.com/GoWxWhoIsTheSpy_$(date "+%Y-%m-%d %H:%M:%S")"
put -r /Users/limengfan/go/src/lmf.mortal.com/GoWxWhoIsTheSpyCopy /home/lighthouse/go/src/lmf.mortal.com
rename /home/lighthouse/go/src/lmf.mortal.com/GoWxWhoIsTheSpyCopy "/home/lighthouse/go/src/lmf.mortal.com/GoWxWhoIsTheSpy"
EOF

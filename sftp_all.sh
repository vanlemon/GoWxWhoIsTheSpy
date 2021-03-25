sftp root@something << EOF
rename /home/lighthouse/go/src/lmf.mortal.com/GoWxWhoIsTheSpy "/home/lighthouse/go/src/lmf.mortal.com/GoWxWhoIsTheSpy_$(date "+%Y-%m-%d %H:%M:%S")"
put -r /Users/limengfan/go/src/lmf.mortal.com/GoWxWhoIsTheSpy /home/lighthouse/go/src/lmf.mortal.com

rename /home/lighthouse/go/src/lmf.mortal.com/GoLogs "/home/lighthouse/go/src/lmf.mortal.com/GoLogs_$(date "+%Y-%m-%d %H:%M:%S")"
put -r /Users/limengfan/go/src/lmf.mortal.com/GoLogs /home/lighthouse/go/src/lmf.mortal.com

rename /home/lighthouse/go/src/lmf.mortal.com/GoLimiter "/home/lighthouse/go/src/lmf.mortal.com/GoLimiter_$(date "+%Y-%m-%d %H:%M:%S")"
put -r /Users/limengfan/go/src/lmf.mortal.com/GoLimiter /home/lighthouse/go/src/lmf.mortal.com
EOF

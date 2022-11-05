echo $show_output_start
((((terraform init -no-color; echo $? >&3) 2>&1 1>/dev/null | appendLogs "/tmp/$s3FileName.txt" >&4) 3>&1) | (read xs; exit $xs)) 4>&1
if [ $? -ne 0 ]; then
  echo $show_output_end
  SaveAndExit "Failed to initialize terraform"
fi
echo $show_output_end
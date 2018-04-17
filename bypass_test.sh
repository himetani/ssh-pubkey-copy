vagrant provision
go run ./*.go bypass bypass -d ./testdata/bypass.json --port 2222 --key $HOME/.ssh/id_rsa.pub -i $HOME/.ssh/id_rsa
go run ./*.go status -d ./testdata/bypass.json --port 2222

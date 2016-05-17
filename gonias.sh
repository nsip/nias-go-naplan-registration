# This is the *nix NIAS batch file launcher. Add extra validators to the bottom of this list. 
# Change the directory as appropriate (go-nias)
# gnatsd MUST be the first program launched

if [ -f "nias.pid" ]
then
echo "There is a nias.pid file in place; run shutdown.sh"
exit
fi

#rem Run the NIAS services. Add to the BOTTOM of this list
# store each PID in pid list
../../nats-io/gnatsd/gnatsd & echo $! > nias.pid

# give the nats server time to come up
sleep 2

cd ./aggregator
./aggregator & echo $! >> ../nias.pid
cd ..


cd ./aslvalidator
./aslvalidator & echo $! >> ../nias.pid
cd ..

cd ./idvalidator
./idvalidator & echo $! >> ../nias.pid
cd ..

cd ./schemavalidator
./schemavalidator & echo $! >> ../nias.pid
cd ..

cd ./dobvalidator
./dobvalidator -tstyr 2016 & echo $! >> ../nias.pid
cd ..

cd ./csvxmlconverter
./csvxmlconverter & echo $! >> ../nias.pid
cd ..

cd ./webui
./webui & echo $! >> ../nias.pid
cd ..

echo "Run the web client (launch browser here):"
echo "http://localhost:8080/nias"


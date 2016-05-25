# This is the *nix NIAS batch file launcher for the dev environment. 
# Add extra validators to the bottom of this list. 
# Change the directory as appropriate (go-nias)
# gnatsd MUST be the first program launched, and is assumed to be on your system

if [ -f "nias.pid" ]
then
echo "There is a nias.pid file in place; run shutdown.sh"
exit
fi

#rem Run the NIAS services. Add to the BOTTOM of this list
# store each PID in pid list
gnatsd & echo $! > nias.pid

# give the nats server time to come up
sleep 2

cd aggregator/
./aggregator & echo $! >> ../nias.pid
cd ../aslvalidator/
./aslvalidator & echo $! >> ../nias.pid
cd ../idvalidator/
./idvalidator & echo $! >> ../nias.pid
cd ../schemavalidator/
./schemavalidator & echo $! >> ../nias.pid
cd ../csvxmlconverter
./csvxmlconverter & echo $! >> ../nias.pid
cd ../dobvalidator/
./dobvalidator -tstyr 2016 & echo $! >> ../nias.pid
cd ../webui/
./webui & echo $! >> ../nias.pid

echo "Run the web client (launch browser here):"
echo "http://localhost:8080/nias"


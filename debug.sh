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

# Build
(cd aggregator; go build)
(cd aslvalidator; go build)
(cd idvalidator; go build)
(cd schemavalidator; go build)
(cd csvxmlconverter; go build)
(cd dobvalidator; go build)
(cd webui; go build)

# Run

(cd aggregator; ./aggregator & echo $! >> ../nias.pid)
./aslvalidator/aslvalidator & echo $! >> nias.pid
./idvalidator/idvalidator & echo $! >> nias.pid
(cd schemavalidator; ./schemavalidator & echo $! >> ../nias.pid)
(cd csvxmlconverter; ./csvxmlconverter & echo $! >> ../nias.pid)
./dobvalidator/dobvalidator -tstyr 2016 & echo $! >> nias.pid
(cd webui; ./webui & echo $! >> ../nias.pid)

echo "Run the web client (launch browser here):"
echo "http://localhost:8080/nias"


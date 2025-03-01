echo "-----"
echo "Checking if server is reachable"
echo "Pinging 5 times"
echo
ping -c 5 http://ec2-52-91-186-123.compute-1.amazonaws.com/

echo "-----"
echo "Checking if API is reachable"
echo
# run curl commands to test that API is working properly
curl http://ec2-52-91-186-123.compute-1.amazonaws.com/

echo "-----"
echo "Checking API predictions with different features"
echo
curl http://ec2-52-91-186-123.compute-1.amazonaws.com/predict?
curl http://ec2-52-91-186-123.compute-1.amazonaws.com/predict?
curl http://ec2-52-91-186-123.compute-1.amazonaws.com/predict?
curl http://ec2-52-91-186-123.compute-1.amazonaws.com/predict?
curl http://ec2-52-91-186-123.compute-1.amazonaws.com/predict?
echo "-----"

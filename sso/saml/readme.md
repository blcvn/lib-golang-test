# 1. Generate certificate for service provider 
cd test/saml
openssl req -x509 -newkey rsa:2048 -keyout myservice.key -out myservice.cert -days 365 -nodes -subj "/CN=myservice.example.com"

cd test/saml
go run main.go 


# 2. Start IDP server 
go run test/saml/idp/idp.go --idp http://localhost:8000

# 3. Register service 
curl localhost:8001/saml/metadata > service_provider.xml 

curl localhost:5001/saml/metadata > service_provider1.xml 

curl -X PUT -T "service_provider.xml" "http://localhost:8000/services/1"
curl -X PUT -T "service_provider.xml" "http://localhost:8000/services/2"

# 2/ You browse to 
localhost:8001/hello

=> The middleware redirects you to http://localhost:8000/sso

email:  alice
password: hunter2

# 3. Redirect to localhost:8001/hello
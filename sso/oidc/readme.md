# Register to hypdra

 hydra create oauth2-client --name "myapp" --redirect-uri http://127.0.0.1:5556/auth/google/callback --token-endpoint-auth-method client_secret_post --grant-type refresh_token,authorization_code --response-type code --scope openid,profile,email  --endpoint http://127.0.0.1:4445

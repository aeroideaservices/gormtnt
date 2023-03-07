docker build -t gormtnt_test_alpine_img .

docker run -dti --env-file=testing/.env --name gormtnt_test_alpine gormtnt_test_alpine_img

go test -v -run TestConnection
docker build --target react-build --tag browsersat:build .
id=$(docker create browsersat:build)
docker cp $id:/react_app/dist build
docker rm -v $id
rm grid/grid -f
cp -f ../../grid grid/
docker build grid -t itfantasy/grid:latest

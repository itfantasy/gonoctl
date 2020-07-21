rm grid/grid-core -f
cp -f ../../grid-core/grid-core grid/
docker build grid -t itfantasy/grid:latest

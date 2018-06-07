#/bin/bash
set -e
[ "$#" -lt 3 ] && echo "Missing args: $0 {SOURCE_PROJECT} {SRC_VERSION} {RELEASE_VERSION}" && exit;

src_project=( $1 )
src_version=( $2 )
release_version=( $3 )

project_src_url="git@github.com:seekplum/"${src_project}".git"

[ $src_project = "plum_exporter" ] && bin_path="client_exporter/plum_exporter" && project="plum_exporter" && dockerimg="quay.io/prometheus/golang-builder"

echo "source code url:"${project_src_url}

root_path=`pwd`/.plum_build && mkdir -p ${root_path} && cd $root_path
export GOPATH=${root_path}
echo "export GOPATH:"${GOPATH}
release_path=/tmp
export_path=${release_path}/$bin_path

mkdir -p ${GOPATH}/src/github.com/seekplum/
mkdir -p ${export_path}


cd ${GOPATH}/src/github.com/seekplum && git clone -b $src_version  $project_src_url && cd $project && git_message=`git log -n 1`

docker run --rm -ti -v $(pwd):/app $dockerimg -i "github.com/seekplum/${project}" -p "linux/amd64"

cp .build/linux-amd64/${src_project} ${export_path}

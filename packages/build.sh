#/bin/bash

RETVAL=0

deploy() {
    home_dir=/home/prometheus
    packages_dir=$home_dir/packages
    curr_path=$(dirname $0)

    # 初始化目录
    rm -rf $home_dir/logs/ && mkdir -p $home_dir/logs
    rm -rf $home_dir/prometheus/ && mkdir -p $home_dir/prometheus/rules
    rm -rf $home_dir/alertmanager/ && mkdir -p $home_dir/alertmanager
    rm -rf $home_dir/grafana/ && mkdir -p $home_dir/grafana/dashboards
    rm -rf $home_dir/node_exporter/ && mkdir -p $home_dir/node_exporter
    rm -rf $home_dir/plum_exporter/ && mkdir -p $home_dir/plum_exporter

    # 解压目录
    tar zxvf $packages_dir/prometheus*.tar.gz --strip-components 1 -C $home_dir/prometheus
    tar zxvf $packages_dir/alertmanager*.tar.gz --strip-components 1 -C $home_dir/alertmanager
    tar zxvf $packages_dir/grafana*.tar.gz --strip-components 1 -C $home_dir/grafana
    tar zxvf $packages_dir/node_exporter*.tar.gz --strip-components 1 -C $home_dir/node_exporter
    tar zxvf $packages_dir/plum_exporter*.tar.gz --strip-components 1 -C $home_dir/plum_exporter

    # 拷贝文件
    scp $curr_path/supervisord.conf /etc/supervisord.conf
    scp $curr_path/prometheus.yml $home_dir/prometheus/
    scp $curr_path/test.rules $home_dir/prometheus/rules/
    scp $curr_path/alertmanager.yml $home_dir/alertmanager/
    scp $curr_path/email.tmpl $home_dir/alertmanager/
    scp $curr_path/defaults.ini $home_dir/grafana/conf/

    supervisorctl shutdown > /dev/null 2&>1
    sleep 5
    supervisord -c /etc/supervisord.conf
}

install() {
    # 安装软件
    yum install -y python-setuptools
    easy_install pip
    pip install supervisor
}

case "$1" in
  deploy)
	deploy
	;;
  install)
	install
	;;
  *)
	deploy  # 默认进行部署
esac

exit $RETVAL

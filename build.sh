#/bin/bash

home_dir=/home/prometheus
packages_dir=$home_dir/packages
mkdir $home_dir
mkdir $packages_dir
mkdir $home_dir/prometheus
mkdir $home_dir/alertmanager
mkdir $home_dir/grafana
mkdir $home_dir/node_exporter
mkdir $home_dir/plum_exporter

scp root@192.168.1.157:/home/prometheus/packages/* $packages_dir

tar zxvf $packages_dir/prometheus*.tar.gz --strip-components 1 -C $home_dir/prometheus
tar zxvf $packages_dir/alertmanager*.tar.gz --strip-components 1 -C $home_dir/alertmanager
tar zxvf $packages_dir/grafana*.tar.gz --strip-components 1 -C $home_dir/grafana
tar zxvf $packages_dir/node_exporter*.tar.gz --strip-components 1 -C $home_dir/node_exporter

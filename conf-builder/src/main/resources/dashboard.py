#-*-coding:utf8-*-
import os
 
#-- dashboard db config --
DASHBOARD_DB_HOST = "${mysql.host}"
DASHBOARD_DB_PORT = ${mysql.port}
DASHBOARD_DB_USER = "${dbuser.dashboard.account}"
DASHBOARD_DB_PASSWD = "${dbuser.password}"
DASHBOARD_DB_NAME = "${dbname.dashboard}"
 
#-- graph db config --
GRAPH_DB_HOST = "${mysql.host}"
GRAPH_DB_PORT = ${mysql.port}
GRAPH_DB_USER = "${dbuser.dashboard.account}"
GRAPH_DB_PASSWD = "${dbuser.password}"
GRAPH_DB_NAME = "${dbname.graph}"
 
#-- portal db config --
PORTAL_DB_HOST = "${mysql.host}"
PORTAL_DB_PORT = ${mysql.port}
PORTAL_DB_USER = "${dbuser.dashboard.account}"
PORTAL_DB_PASSWD = "${dbuser.password}"
PORTAL_DB_NAME = "${dbname.portal}"

#-- app config --
DEBUG = ${log.debug.python}
SECRET_KEY = "2mf09vjRDC"
SESSION_COOKIE_NAME = "open-falcon"
PERMANENT_SESSION_LIFETIME = 3600 * 24 * 30
SITE_COOKIE = "open-falcon-ck"
 
#-- query config --
QUERY_ADDR = "${url.query}"
 
BASE_DIR = "${path.dashboard.base}"
LOG_PATH = os.path.join(BASE_DIR,"log/")
 
JSONCFG = {}
JSONCFG['database'] = {}
JSONCFG['database']['host']     = '${mysql.host}'
JSONCFG['database']['port']     = '${mysql.port}'
JSONCFG['database']['account']  = '${dbuser.dashboard.account}'
JSONCFG['database']['password'] = '${dbuser.password}'
JSONCFG['database']['db']       = '${dbname.uic}'
JSONCFG['database']['table']    = 'session'

JSONCFG['shortcut'] = {}
JSONCFG['shortcut']['falconPortal']     = "${url.portal}"
JSONCFG['shortcut']['falconDashboard']  = "${url.dashboard}"
JSONCFG['shortcut']['grafanaDashboard'] = "${url.grafana}"
JSONCFG['shortcut']['falconAlarm']      = "${url.alarm}"
JSONCFG['shortcut']['falconUIC']        = "${url.fe}"

JSONCFG['redirectUrl'] = '${url.fe}/auth/login?callback=${url.dashboard}/'

try:
    from rrd.local_config import *
except:
    pass

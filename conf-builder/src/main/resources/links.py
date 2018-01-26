# -*- coding:utf-8 -*-

# -- app config --
DEBUG = ${log.debug.python}

# -- db config --
DB_HOST = "${mysql.host}"
DB_PORT = ${mysql.port}
DB_USER = "${dbuser.links.account}"
DB_PASS = "${dbuser.password}"
DB_NAME = "${dbname.links}"

# -- cookie config --
SECRET_KEY = "mfiovn2FfA1yhb"
SESSION_COOKIE_NAME = "falcon-links"
PERMANENT_SESSION_LIFETIME = 3600 * 24 * 30

try:
    from frame.local_config import *
except Exception, e:
    print 'level=warning msg="%s"' % e

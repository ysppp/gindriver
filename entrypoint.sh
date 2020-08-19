#!/bin/sh

# Debug mode
if [ "$DEBUG" ]; then
    set -x
fi

# Time interval
if [ -z "$RESTART" ]; then
    RESTART=60
fi

# Gin mode
if [ -z "$GIN_MODE" ]; then
    GIN_MODE=release
fi

# Test flag
if [ -z "$FLAG" ]; then
    FLAG="flag{Th1s_is_an_3xamp1e_fl114g}"
fi

if [ -z "$SFLAG" ]; then
    SFLAG="flag{Th1s_is_an0th3r_3xamp1e_fl114g}"
fi

echo $FLAG > /flag
echo "The End?" >> /flag

# Secret flag
SFLAGPATH="/etc/secret$(cat /dev/urandom | head -n 10 | md5sum | head -c 16)flag"
echo $SFLAG > "$SFLAGPATH"


# Init RP config
if [ -z "$RPID" ]; then
    RPID="gindriver.evi0s.com"
fi

if [ -z "$RPORIGIN" ]; then
    RPORIGIN="https://gindriver.evi0s.com/"
fi

if [ -z "$RPNAME" ]; then
    RPNAME="GinDriver"
fi

# Start services
/etc/init.d/mysql start
/etc/init.d/ssh start

# Init database
if [ -z "$DBNAME" ]; then
    DBNAME=gindriver
fi

if [ -z "$DBUSER" ]; then
    DBUSER=ctf
fi

if [ -z "$DBPASS" ]; then
    DBPASS=$(cat /dev/urandom | head -n 10 | md5sum | head -c 32)
    echo "Database init password: ${DBPASS}"
fi

# Create database users
mysql -e "CREATE DATABASE ${DBNAME};"
mysql -e "GRANT USAGE ON *.* TO '${DBUSER}'@'localhost' IDENTIFIED BY '${DBPASS}' WITH GRANT OPTION; GRANT SELECT,INSERT,UPDATE,DELETE,CREATE,DROP ON ${DBNAME}.* TO '${DBUSER}'@'localhost' IDENTIFIED BY '${DBPASS}'; GRANT EXECUTE ON ${DBNAME}.* TO '${DBUSER}'@'localhost' IDENTIFIED BY '${DBPASS}'; FLUSH PRIVILEGES;"

cat > /backup/config/config.yml << EOF
appname: gindriver

db:
  name: $DBNAME
  user: $DBUSER
  pass: $DBPASS
  host: "127.0.0.1"
  port: 3306
  param: "charset=utf8&loc=Local"

listenaddr: "0.0.0.0:3000"
rpid: "$RPID"
rporigin: "$RPORIGIN"
rpdisplayname: "$RPNAME"

EOF

# Admin public key
mkdir -p /backup/public/pubkeys
echo "A_sup3r_Sup3r_S3cret_PUBBBBBBB_k3y?" > /backup/public/pubkeys/admin.pub

while true; do
    # Backup public keys
    cp -r /var/www/public/pubkeys/* /backup/public/pubkeys/

    # Clean up
    rm -rf /var/www/*

    # Add hint
    mkdir -p /var/www/.ssh
    echo "Interesting, right?" > /var/www/.ssh/authorized_keys

    # Copy to dist
    cp -r /backup/* /var/www


    # Fix ownership & permission
    chown -R www-data.www-data /var/www/
    chmod -R 0755 /var/www
    chmod +x /var/www/gindriver

    # Run
    su - www-data -s "/bin/sh" -c "GIN_MODE=$GIN_MODE /var/www/gindriver" &

    # Time interval
    sleep $RESTART

    # Time to die!
    killall gindriver

    # Take a breath
    sleep 5
done


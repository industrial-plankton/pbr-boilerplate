# From Scratch Droplet Setup

## User Setup

Login

    ssh root@your_server_ip

Add user

    adduser USERNAME
    enter new password: PASS

Give admin privlege

    usermod -aG sudo USERNAME

## Firewall

    ufw LIMIT OpenSSH comment "SSH connection, port 22"
    ufw ALLOW https comment "Main website entrypoint, port 443"
    ufw LIMIT http comment "Redirects to HTTPS, port 80"
    ufw enable

HTTPS could be LIMIT'd but the tolerance would need to be increased, *I easily hit it while testing*

## Enable ssh for new user

    rsync --archive --chown=USERNAME:USERNAME ~/.ssh /home/USERNAME

## Update packges for postgres

    sudo apt update && sudo apt -y upgrade
    wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | sudo apt-key add -
    RELEASE=$(lsb_release -cs)
    echo "deb http://apt.postgresql.org/pub/repos/apt/ ${RELEASE}"-pgdg main | sudo tee  /etc/apt/sources.list.d/pgdg.list
    sudo apt update
    sudo apt -y install postgresql-11

## To allow network access, edit configuration file: **not nessesary**

    sudo nano /etc/postgresql/11/main/postgresql.conf

Add below to CONNECTIONS AND AUTHENTICATION section.

    listen_addresses = '*'

You can also specify server IP Address

    listen_addresses = '192.168.17.12'
restart service after making a change

    sudo systemctl restart postgresql
If you have an active UFW firewall, allow port 5432

    sudo ufw allow 5432/tcp

## Set databse password for admin user

    sudo su - postgres
    psql -c "alter user postgres with password 'StrongPassword'"

## Backup/Restore

### From postgre computer account

    pg_dump -U DATABASEUSER DATABASENAME > BACKUPFILE

### Automaticly with cron job

google it, i didnt write it down

### To use the backup file

    createdb -T template0 NEWDATABASENAME
    psql -d NEWDATABASENAME -1 -f BACKUPFILE

### Continuous backups

open postgresql.conf and set:

    wal_level = replica
    archive_mode = on
    archive_command = 'test ! -f /mnt/backup/%f && cp %p /mnt/backup/%f'

then

    chmod /mnt/backup/ 
so the postgres user has write access, i just gave everyone free access with:

    sudo chmod -R 777 /mnt/backup/

## Database side changelog use

[Follow this guide](https://wiki.postgresql.org/wiki/Audit_trigger_91plus)

## HTTPS Certificate

    sudo apt-get update
    sudo apt-get install software-properties-common
    sudo add-apt-repository universe
    sudo add-apt-repository ppa:certbot/certbot
    sudo apt-get update

    sudo apt-get install certbot
    sudo certbot certonly --standalone

## Go_Server

put `credentials.json` and `connection.md` in same folder that backend will run

`credentials.json` comes form google cloud platform, and is what allows authorization to use the sheets API alongside the Oauth2 tokens.

`connection.md` just contains the Database login info

## Run Go_server as a service

    sudo nano /etc/systemd/system/Go_Server.service     

Paste:

    [Unit]
    Description=Golang Web Server for Database interaction service
    After=network.target
    After=postgresql.service
    StartLimitIntervalSec=0

    [Service]
    Type=simple
    Restart=Always
    RestartSec=2
    ExecStart=/home/cameron/Go_Server/backend

    [Install]
    WantedBy=multi-user.target

Then

    sudo systemctl enable Go_Server

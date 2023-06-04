export USER_NAME=user_name
export USER_PWD=user_password

echo "Creating user with name $USER_NAME"
adduser --quiet --disabled-password --shell /bin/bash --home /home/$USER_NAME --gecos "$USER_NAME" $USER_NAME
echo  "$USER_NAME:$USER_PWD" | chpasswd

echo "User $USER_NAME created"

echo "Setting up ssh"
sed -i 's/#PubkeyAuthentication yes/PubkeyAuthentication yes/g' /etc/ssh/sshd_config
systemctl restart ssh

echo "Installing dependencies"
apt update
# Y, 2
#apt upgrade -y
apt install curl sudo -y

echo "Installing docker"
curl -fsSL https://get.docker.com -o get-docker.sh
sh get-docker.sh
rm get-docker.sh
usermod -aG docker $USER_NAME
echo "Docker installed"
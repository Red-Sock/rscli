export CONNECT_NAME=connection_name
export USER_NAME=user_name
export HOST=host_api
export PORT=port_number

rm -f ~/.ssh/$CONNECT_NAME*
ssh-keygen -t rsa -f ~/.ssh/$CONNECT_NAME -q -N ""
ssh-copy-id  -p $PORT -i ~/.ssh/$CONNECT_NAME.pub $USER_NAME@$HOST

echo "
Host $CONNECT_NAME $HOST
  HostName $HOST
  IdentityFile ~/.ssh/$CONNECT_NAME
  User $USER_NAME
  Port $PORT
" >> ./config


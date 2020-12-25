### TODO ADD MORE BELOW kubectx, kubectl, helm3, ... ? ###

echo "installing kubectx........."

if [[ "$OSTYPE" == "linux-gnu"* ]]; then
  # Adds the stuff you needd to a linux node...
  wget https://github.com/ahmetb/kubectx/releases/download/v0.9.1/kubectx_v0.9.1_linux_x86_64.tar.gz
  tar -xvf kubectx_v0.9.1_linux_x86_64.tar.gz
  cp kubectx /usr/local/bin/
  sudo cp kubectx /usr/local/bin/
elif [[ "$OSTYPE" == "darwin"* ]]; then
  # Adds the stuff you needd to a linux node...
  wget https://github.com/ahmetb/kubectx/releases/download/v0.9.1/kubectx_v0.9.1_darwin_x86_64.tar.gz
  tar -xvf kubectx_v0.9.1_linux_x86_64.tar.gz
  cp kubectx /usr/local/bin/
  sudo cp kubectx /usr/local/bin/
fi

kubectx

echo "................... done!"

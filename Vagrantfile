# The resources to give the VM.
CPUS = 2
MEMORY = 4096

# The name of the VM.
NAME = "violet-test"

Vagrant.configure("2") do |config|
    config.vm.box = "generic/alpine317"

    config.vm.define NAME
    config.vm.hostname = NAME

    config.vm.synced_folder ".", "/vagrant", disabled: false, type: "virtiofs"

    config.vm.provision "setup", type: "shell", keep_color: true, inline: <<-SCRIPT
        # Set up shell profile
        cp /vagrant/test/bashrc /etc/profile.d/bashrc.sh

        # Update packages
        apk update

        # Fetch Gum
        wget -qO- \
            https://github.com/charmbracelet/gum/releases/download/v0.10.0/gum_0.10.0_Linux_x86_64.tar.gz \
            | tar -xz -C /usr/bin/ gum

        # Install and start Docker
        apk add --no-cache docker
        rc-update add docker default
        rc-service docker start
        addgroup vagrant docker

        # Install Vagrant from source
        apk add --no-cache ruby ruby-dev ruby-bundler git gcc make libc-dev ncurses go
        git clone https://github.com/hashicorp/vagrant.git
        cd vagrant
            bundle install
            bundle --binstubs exec
            ln -sf $(pwd)/exec/vagrant /usr/bin/vagrant

            vagrant version
    SCRIPT
    config.vm.provision "install", type: "shell", keep_color: true, privileged: false, inline: <<-SCRIPT
        cd /vagrant
            make install
    SCRIPT
    config.vm.provider :libvirt do |l, override|
        l.driver = "kvm"
        l.cpus = CPUS
        l.memory = MEMORY
        l.disk_bus = "virtio"
        l.qemu_use_session = false

        l.default_prefix = ""

        l.nested = true

        l.memorybacking :access, :mode => "shared"

        # Enable Hyper-V enlightments: https://blog.wikichoon.com/2014/07/enabling-hyper-v-enlightenments-with-kvm.html
        l.hyperv_feature :name => 'relaxed', :state => 'on'
        l.hyperv_feature :name => 'synic',   :state => 'on'
        l.hyperv_feature :name => 'vapic',   :state => 'on'
        l.hyperv_feature :name => 'vpindex', :state => 'on'
      end
end
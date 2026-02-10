complete -x -c smug -a "(ls ~/.config/smug | grep -v \"smug\.log\" | sed -e 's/\..*//')"
complete -c smug -n '__fish_use_subcommand' -a 'rm' -d 'Remove project configuration'
complete -c smug -n '__fish_use_subcommand' -a 'switch' -d 'Switch to a project session'

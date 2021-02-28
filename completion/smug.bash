_smug_list_projects () {
  if [ "${#COMP_WORDS[@]}" != "2" ]; then
    return
  fi
  COMPREPLY=($(ls ~/.config/smug | grep -v "smug\.log" | sed -e 's/\..*//'))
}

complete -F _smug_list_projects smug

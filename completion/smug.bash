#! /usr/bin/env bash
_smug() {
    local ISF=$'\n'
    local reply

    local cur="${COMP_WORDS[COMP_CWORD]}"
    local prev="${COMP_WORDS[COMP_CWORD-1]}"

    # if command is 'list' or 'print' do not suggest more
    for word in ${COMP_WORDS[@]}; do
        case $word in
            list|print|rm) return
        esac
    done

    # commands
    if (( "${#COMP_WORDS[@]}" == 2 )); then
        reply=($(compgen -W "list print rm start stop switch" -- "${cur}"))
    fi

    # projects
    if (( "${#COMP_WORDS[@]}" == 3 )); then
        case ${prev} in
            start|stop|rm|switch)
                reply=($(compgen -W "$(smug list | grep -F -v smug)" -- "${cur}"))
        esac
    fi

    # options
    if (( "${#COMP_WORDS[@]}" > 3 )); then
        local options=( "--file" "--windows" "--attach" "--debug" )

        # --windows waits for a list
        case $prev in
            -w|--windows) return
        esac

        # suggest options that were not specified already
        for word in "${COMP_WORDS[@]}"; do
            case $word in
                -f|--file) options=( "${options[@]/--file}" ) ;;
                -w|--windows) options=( "${options[@]/--windows}" ) ;;
                -a|--attach) options=( "${options[@]/--attach}" ) ;;
                -d|--debug) options=( "${options[@]/--debug}" ) ;;
            esac
        done

        # array to string
        local options="$(echo ${options[@]})"

        reply=($(compgen -W "${options}" -- "${cur}"))
    fi


    # if only one match proceed with autocompletion
    if (( "${#reply[@]}" == 1 )); then
        COMPREPLY=( "${reply[0]}" )
    else
        # when 'TAB TAB' is pressed
        if (( COMP_TYPE == 63 )); then
            # print suggestions as list with padding
            for i in "${!reply[@]}"; do
                reply[$i]="$(printf '%*s' "-$COLUMNS" "${reply[$i]}")"
            done
        fi
        # print suggestions
        COMPREPLY=( "${reply[@]}" )
    fi
}

complete -F _smug smug

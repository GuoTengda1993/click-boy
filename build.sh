build_mac() {
    wails build
    rm -f click-boy.dmg
    appdmg dmg.json click-boy.dmg
    mv click-boy.dmg build/bin
}

build_win() {
    wails build --platform windows/amd64 -skipbindings

build() {
    echo "~~~ building mac start ~~~"
    build_mac
    echo "*** building mac finish ***"
    echo ""

    echo "~~~ building win start ~~~"
    build_win
    echo "*** building win finish ***"
    echo ""
}

case "$1" in
mac)
    build_mac
    ;;
win)
    build_win
    ;;
*)
    build
    ;;
esac

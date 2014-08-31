package config

const (
    APP_NAME = "Tindra"
    TAGLINE  = "the glimmering site generator!"
    VERSION  = "0.2"

    MAIN_CONFIG_FILENAME = "config.yml"

    // input dirs
    INCLUDES_DIR_NAME = "includes"
    LAYOUTS_DIR_NAME  = "layouts"
    POSTS_DIR_NAME    = "posts"
    DATA_DIR_NAME     = "data"
    PLUGINS_DIR_NAME  = "plugins"
    CSS_DIR_NAME      = "css"
    JS_DIR_NAME       = "js"
    RESOURCE_DIR_NAME = "media"

    // output dirs
    BUILD_DIR_NAME = "build"
)

const (
    INIT      = 1 << iota
    RAW       = 1 << iota
    BUILT     = 1 << iota
    INCLUDES  = 1 << iota
    GENERATED = 1 << iota
)

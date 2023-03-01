project('yggdrasil-worker-package-manager', version: '0.2.0')

go = find_program('go')

dbus = dependency('dbus-1')

config_data = configuration_data({
  'libexecdir': join_paths(get_option('prefix'), get_option('libexecdir')),
  'configdir': join_paths(get_option('sysconfdir'), meson.project_name()),
})

goldflags = ''
goldflags += ' -X "main.Version=' + meson.project_version() + '"'
goldflags += ' -X "main.ConfigDir=' + config_data.get('configdir') + '"'

gobuildflags = [
  '-buildmode=pie',
]

custom_target(meson.project_name(),
  build_always_stale: true,
  output: meson.project_name(),
  command: [go, 'build', gobuildflags, '-o', '@OUTPUT@', '-ldflags', goldflags, 'github.com/redhatinsights/yggdrasil-worker-package-manager'],
  install: true,
  install_dir: join_paths(get_option('libexecdir'))
)

configure_file(
  configuration: config_data,
  input: 'com.redhat.yggdrasil.Worker1.package_manager.service.in',
  output: 'com.redhat.yggdrasil.Worker1.package_manager.service',
  install: true,
  install_dir: dbus.get_pkgconfig_variable('system_bus_services_dir')
)

install_data('config.toml',
  install_dir: config_data.get('configdir')
)

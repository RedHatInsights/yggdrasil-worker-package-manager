project('yggdrasil-worker-package-manager', version: '0.2.2')

go = find_program('go')

dbus = dependency('dbus-1')

if get_option('vendor')
  meson.add_dist_script(join_paths('build-aux', 'vendor.sh'))
endif

config_data = configuration_data({
  'libexecdir': join_paths(get_option('prefix'), get_option('libexecdir')),
  'configdir': join_paths(get_option('sysconfdir'), meson.project_name()),
})

goldflags = get_option('goldflags')
goldflags += ' -X "main.Version=' + meson.project_version() + '"'
goldflags += ' -X "main.ConfigDir=' + config_data.get('configdir') + '"'

gobuildflags = get_option('gobuildflags')

custom_target(meson.project_name(),
  build_always_stale: true,
  output: meson.project_name(),
  command: [go, 'build', gobuildflags, '-o', '@OUTPUT@', '-ldflags', goldflags, 'github.com/redhatinsights/yggdrasil-worker-package-manager'],
  install: true,
  install_dir: join_paths(get_option('libexecdir'))
)

configure_file(
  configuration: config_data,
  input: 'com.redhat.Yggdrasil1.Worker1.package_manager.service.in',
  output: 'com.redhat.Yggdrasil1.Worker1.package_manager.service',
  install: true,
  install_dir: dbus.get_pkgconfig_variable('system_bus_services_dir')
)

install_data('config.toml',
  install_dir: config_data.get('configdir')
)

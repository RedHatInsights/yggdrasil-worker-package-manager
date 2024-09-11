project('yggdrasil-worker-package-manager', version: '0.2.3')

go = find_program('go')

dbus = dependency('dbus-1')
systemd = dependency('systemd')

if get_option('vendor')
  meson.add_dist_script(join_paths('build-aux', 'vendor.sh'))
endif

config_data = configuration_data(
  {
    'libexecdir': join_paths(get_option('prefix'), get_option('libexecdir')),
    'configdir': join_paths(get_option('sysconfdir'), meson.project_name()),
    'worker_user': get_option('worker_user'),
  },
)

goldflags = get_option('goldflags')
goldflags += ' -X "main.Version=' + meson.project_version() + '"'
goldflags += ' -X "main.ConfigDir=' + config_data.get('configdir') + '"'

gobuildflags = get_option('gobuildflags')

custom_target(
  meson.project_name(),
  build_always_stale: true,
  output: meson.project_name(),
  command: [
    go,
    'build',
    gobuildflags,
    '-o', '@OUTPUT@',
    '-ldflags', goldflags,
    'github.com/redhatinsights/yggdrasil-worker-package-manager',
  ],
  install: true,
  install_dir: join_paths(get_option('libexecdir')),
)

install_data('config.toml', install_dir: config_data.get('configdir'))

subdir('data')
subdir('dist')

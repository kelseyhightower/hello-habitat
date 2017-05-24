pkg_name="hello-habitat"
pkg_description="Hello Habitat demo application."
pkg_origin="kelseyhightower"
pkg_version="0.1.0"
pkg_license=('Apache-2.0')
pkg_upstream_url="https://github.com/kelseyhightower/hello-habitat"

pkg_svc_run="hello-habitat --config-file ${pkg_svc_config_path}/config.json"
pkg_svc_user="hab"
pkg_svc_group="${pkg_svc_user}"
pkg_bin_dirs=(bin)

pkg_exports=(
  [http]=server.http
)
pkg_exposes=(http)

pkg_build_deps=(core/go)

do_build() {
  pushd "${PLAN_CONTEXT}"/.. > /dev/null
  go build -o hello-habitat . 
  popd > /dev/null
}

do_install() {
  mkdir -p "${pkg_prefix}/bin"
  cp "${PLAN_CONTEXT}"/../hello-habitat "${pkg_prefix}/bin"
}

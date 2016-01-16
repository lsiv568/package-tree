class PackageRepository
  def initialize
    @packages = {}
  end

  def execute(cmd)
    method_name = cmd.command.to_s.downcase.to_sym
    send(method_name, cmd)
  end

  private
  def install(cmd)
    if (cmd.dependencies - all_installed_packages).empty? then
      @packages[cmd.package] = cmd.dependencies
      true
    else
      false
    end
  end

  def uninstall(cmd)
    if ! any_package_depends_on?(cmd.package) then
      @packages.delete(cmd.package)
      true
    else
      false
    end
  end

  def query(cmd)
    @packages.has_key?(cmd.package)
  end

  def all_installed_packages
    @packages.keys
  end

  def any_package_depends_on?(pkg)
    @packages.values.flatten.include?(pkg)
  end
end

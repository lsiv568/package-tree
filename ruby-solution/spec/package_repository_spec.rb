require 'rspec'
require 'command'
require 'package_repository'

describe PackageRepository do
  let (:package_repository) { PackageRepository.new }
  let(:first_install) { Command.new('INSTALL|banana-tree|') }
  let(:first_query) { Command.new('QUERY|banana-tree|') }
  let(:first_uninstall) { Command.new('UNINSTALL|banana-tree|') }
  let(:second_install) { Command.new('INSTALL|banana|banana-tree') }
  let(:second_query) { Command.new('QUERY|banana|') }
  let(:second_uninstall) { Command.new('UNINSTALL|banana|') }

  it 'installs a package if dependencies are met' do
    expect(package_repository.execute(first_install)).to eq(true)
    expect(package_repository.execute(first_query)).to eq(true)

    expect(package_repository.execute(second_install)).to eq(true)
    expect(package_repository.execute(second_query)).to eq(true)
  end

  it 'does not install a package if dependencies arent met' do
    expect(package_repository.execute(second_install)).to eq(false)
  end

  it 'uninstalls existing package with no inbound dependencies' do
    package_repository.execute(first_install)
    package_repository.execute(second_install)

    expect(package_repository.execute(second_uninstall)).to eq(true)
    expect(package_repository.execute(second_query)).to eq(false)

    expect(package_repository.execute(first_uninstall)).to eq(true)
    expect(package_repository.execute(first_query)).to eq(false)
  end

  it 'doesnt uninstall existing package with inbound dependencies' do
    package_repository.execute(first_install)
    package_repository.execute(second_install)

    expect(package_repository.execute(first_uninstall)).to eq(false)
    expect(package_repository.execute(first_query)).to eq(true)
  end
end

dao_path="labsystem/dao"
srv_path="labsystem/service"

test_dao="go test "$dao_path
test_srv="go test "$srv_path

# -------service---------
$test_srv/admin
$test_srv/user
$test_srv/class

# -------dao---------
$test_dao/admin
$test_dao/user
$test_dao/class
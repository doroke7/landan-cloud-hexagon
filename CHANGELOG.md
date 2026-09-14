1. 把 proto/resource/logic/admin_permission_group.proto 裡面的 
AdminPermissionGroupValue
AdminPermissionValue
移動到  proto/resource/common.proto


2. 也把 proto/resource/logic/admin_user.proto 裡面的
AdminUserValue 
也移動 到 proto/resource/common.proto

3. 其他 proto/resource/logic/app_user.proto
proto/resource/event/admin_user.proto
的 XXXXValue 也移動到  proto/resource/common.proto， 其中 XXXX 是名字

4. proto/resource/model 底下的 XXXXValue 也移動到 proto/resource/common.proto
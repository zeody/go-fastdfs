package go_fastdfs

const (
	FDFS_PROTO_CMD_QUIT                                     byte = 82
	TRACKER_PROTO_CMD_SERVER_LIST_GROUP                     byte = 91
	TRACKER_PROTO_CMD_SERVER_LIST_STORAGE                   byte = 92
	TRACKER_PROTO_CMD_SERVER_DELETE_STORAGE                 byte = 93
	TRACKER_PROTO_CMD_SERVICE_QUERY_STORE_WITHOUT_GROUP_ONE byte = 101
	TRACKER_PROTO_CMD_SERVICE_QUERY_FETCH_ONE               byte = 102
	TRACKER_PROTO_CMD_SERVICE_QUERY_UPDATE                  byte = 103
	TRACKER_PROTO_CMD_SERVICE_QUERY_STORE_WITH_GROUP_ONE    byte = 104
	TRACKER_PROTO_CMD_SERVICE_QUERY_FETCH_ALL               byte = 105
	TRACKER_PROTO_CMD_SERVICE_QUERY_STORE_WITHOUT_GROUP_ALL byte = 106
	TRACKER_PROTO_CMD_SERVICE_QUERY_STORE_WITH_GROUP_ALL    byte = 107
	TRACKER_PROTO_CMD_RESP                                  byte = 100
	FDFS_PROTO_CMD_ACTIVE_TEST                              byte = 111
	STORAGE_PROTO_CMD_UPLOAD_FILE                           byte = 11
	STORAGE_PROTO_CMD_DELETE_FILE                           byte = 12
	STORAGE_PROTO_CMD_SET_METADATA                          byte = 13
	STORAGE_PROTO_CMD_DOWNLOAD_FILE                         byte = 14
	STORAGE_PROTO_CMD_GET_METADATA                          byte = 15
	STORAGE_PROTO_CMD_UPLOAD_SLAVE_FILE                     byte = 21
	STORAGE_PROTO_CMD_QUERY_FILE_INFO                       byte = 22
	STORAGE_PROTO_CMD_UPLOAD_APPENDER_FILE                  byte = 23
	STORAGE_PROTO_CMD_APPEND_FILE                           byte = 24
	STORAGE_PROTO_CMD_MODIFY_FILE                           byte = 34
	STORAGE_PROTO_CMD_TRUNCATE_FILE                         byte = 36
	STORAGE_PROTO_CMD_RESP                                  byte = TRACKER_PROTO_CMD_RESP
	FDFS_STORAGE_STATUS_INIT                                byte = 0
	FDFS_STORAGE_STATUS_WAIT_SYNC                           byte = 1
	FDFS_STORAGE_STATUS_SYNCING                             byte = 2
	FDFS_STORAGE_STATUS_IP_CHANGED                          byte = 3
	FDFS_STORAGE_STATUS_DELETED                             byte = 4
	FDFS_STORAGE_STATUS_OFFLINE                             byte = 5
	FDFS_STORAGE_STATUS_ONLINE                              byte = 6
	FDFS_STORAGE_STATUS_ACTIVE                              byte = 7
	FDFS_STORAGE_STATUS_NONE                                byte = 99

	STORAGE_SET_METADATA_FLAG_OVERWRITE byte = 'O'

	STORAGE_SET_METADATA_FLAG_MERGE      byte   = 'M'
	FDFS_PROTO_PKG_LEN_SIZE              int64  = 8
	FDFS_PROTO_HEADER_LEN                int64  = FDFS_PROTO_PKG_LEN_SIZE + 2
	FDFS_PROTO_CMD_SIZE                  int    = 1
	FDFS_GROUP_NAME_MAX_LEN              int64  = 16
	FDFS_IPADDR_SIZE                     int64  = 16
	FDFS_DOMAIN_NAME_MAX_SIZE            int64  = 128
	FDFS_VERSION_SIZE                    int64  = 6
	FDFS_STORAGE_ID_MAX_SIZE             int64  = 16
	FDFS_RECORD_SEPERATOR                string = "\u0001"
	FDFS_FIELD_SEPERATOR                 string = "\u0002"
	TRACKER_QUERY_STORAGE_FETCH_BODY_LEN        = FDFS_GROUP_NAME_MAX_LEN + FDFS_IPADDR_SIZE - 1 + FDFS_PROTO_PKG_LEN_SIZE
	TRACKER_QUERY_STORAGE_STORE_BODY_LEN        = FDFS_GROUP_NAME_MAX_LEN + FDFS_IPADDR_SIZE + FDFS_PROTO_PKG_LEN_SIZE
	FDFS_FILE_EXT_NAME_MAX_LEN           byte   = 6
	FDFS_FILE_PREFIX_MAX_LEN             byte   = 16
	FDFS_FILE_PATH_LEN                   byte   = 10
	FDFS_FILENAME_BASE64_LENGTH          byte   = 27
	FDFS_TRUNK_FILE_INFO_LEN             byte   = 16
	ERR_NO_ENOENT                        byte   = 2
	ERR_NO_EIO                           byte   = 5
	ERR_NO_EBUSY                         byte   = 16
	ERR_NO_EINVAL                        byte   = 22
	ERR_NO_ENOSPC                        byte   = 28
	ECONNREFUSED                         byte   = 61
	ERR_NO_EALREADY                      byte   = 114
	INFINITE_FILE_SIZE                   int64  = 256 * 1024 * 1024 * 1024 * 1024 * 1024
	APPENDER_FILE_SIZE                          = INFINITE_FILE_SIZE
	TRUNK_FILE_MARK_SIZE                 int64  = 512 * 1024 * 1024 * 1024 * 1024 * 1024
	NORMAL_LOGIC_FILENAME_LENGTH                = FDFS_FILE_PATH_LEN + FDFS_FILENAME_BASE64_LENGTH + FDFS_FILE_EXT_NAME_MAX_LEN + 1
	TRUNK_LOGIC_FILENAME_LENGTH                 = NORMAL_LOGIC_FILENAME_LENGTH + FDFS_TRUNK_FILE_INFO_LEN
	PROTO_HEADER_CMD_INDEX                      = FDFS_PROTO_PKG_LEN_SIZE
	PROTO_HEADER_STATUS_INDEX                   = FDFS_PROTO_PKG_LEN_SIZE + 1
)

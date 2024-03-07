package pkg

const (
	OS_READ        = 04
	OS_WRITE       = 02
	OS_EX          = 01
	OS_USER_SHIFT  = 6
	OS_GROUP_SHIFT = 3
	OS_OTH_SHIFT   = 0

	OS_USER_R   = OS_READ << OS_USER_SHIFT  // User read permission
	OS_USER_W   = OS_WRITE << OS_USER_SHIFT // User write permission
	OS_USER_X   = OS_EX << OS_USER_SHIFT    // User execute permission
	OS_USER_RW  = OS_USER_R | OS_USER_W     // User read-write permission
	OS_USER_RWX = OS_USER_RW | OS_USER_X    // User read-write-execute permission

	OS_GROUP_R   = OS_READ << OS_GROUP_SHIFT  // Group read permission
	OS_GROUP_W   = OS_WRITE << OS_GROUP_SHIFT // Group write permission
	OS_GROUP_X   = OS_EX << OS_GROUP_SHIFT    // Group execute permission
	OS_GROUP_RW  = OS_GROUP_R | OS_GROUP_W    // Group read-write permission
	OS_GROUP_RWX = OS_GROUP_RW | OS_GROUP_X   // Group read-write-execute permission

	OS_OTH_R   = OS_READ << OS_OTH_SHIFT  // Others read permission
	OS_OTH_W   = OS_WRITE << OS_OTH_SHIFT // Others write permission
	OS_OTH_X   = OS_EX << OS_OTH_SHIFT    // Others execute permission
	OS_OTH_RW  = OS_OTH_R | OS_OTH_W      // Others read-write permission
	OS_OTH_RWX = OS_OTH_RW | OS_OTH_X     // Others read-write-execute permission

	OS_ALL_R   = OS_USER_R | OS_GROUP_R | OS_OTH_R    // All users read permission
	OS_ALL_W   = OS_USER_W | OS_GROUP_W | OS_OTH_W    // All users write permission
	OS_ALL_X   = OS_USER_X | OS_GROUP_X | OS_OTH_X    // All users execute permission
	OS_ALL_RW  = OS_USER_RW | OS_GROUP_RW | OS_OTH_RW // All users read-write permission
	OS_ALL_RWX = OS_ALL_RW | OS_GROUP_X               // All users read-write-execute permission
)

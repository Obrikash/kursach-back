DELETE FROM user_groups
WHERE (user_id = 1 AND group_id = 1)
   OR (user_id = 2 AND group_id = 2)
   OR (user_id = 3 AND group_id = 3)
   OR (user_id = 4 AND group_id = 4)
   OR (user_id = 1 AND group_id = 3)
   OR (user_id = 2 AND group_id = 4);
DELETE FROM training_groups
WHERE (pool_id = 1 AND category_id = 1 AND trainer_id = 1)
   OR (pool_id = 1 AND category_id = 3 AND trainer_id = 2)
   OR (pool_id = 2 AND category_id = 2 AND trainer_id = 3)
   OR (pool_id = 2 AND category_id = 4 AND trainer_id = 4);
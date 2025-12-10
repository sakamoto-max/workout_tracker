-- CREATE TABLE USERS (
-- 	ID SERIAL PRIMARY KEY,
-- 	NAME TEXT NOT NULL,
-- 	EMAIL TEXT NOT NULL,
-- 	HASHED_PASSWORD TEXT NOT NULL,
-- 	ROLE TEXT NOT NULL,
-- 	CREATED_AT TIMESTAMPTZ,
-- 	UPDATED_AT TIMESTAMPTZ
-- )


-- CREATE TABLE EXERCISES (
-- 	ID SERIAL PRIMARY KEY,
-- 	EXERCISE_NAME TEXT NOT NULL,
-- 	TYPE TEXT NOT NULL,
-- 	BODY_PART TEXT NOT NULL
-- )

-- CREATE TABLE WORKOUTPLANS(
-- 	USER_ID INTEGER REFERENCES USERS(ID),
-- 	PLAN_ID INTEGER REFERENCES PLAN(ID),
-- 	PLAN_NAME TEXT	
-- )

-- CREATE TABLE PLAN(
-- 	ID SERIAL PRIMARY KEY,
-- 	EXERCISE_NAME TEXT NOT NULL,
-- 	EXERCISE_TRACKER_ID INT REFERENCES EXERCISE_TRACKER(ID)
-- )

-- CREATE TABLE EXERCISE_TRACKER(
-- 	ID SERIAL PRIMARY KEY,
-- 	SET_1_REPS INTEGER,
-- 	WEIGHT_FOR_SET_1 INTEGER,
-- 	SET_2_REPS INTEGER,
-- 	WEIGHT_FOR_SET_2 INTEGER,
-- 	SET_3_REPS INTEGER,
-- 	WEIGHT_FOR_SET_3 INTEGER,
-- 	SET_4_REPS INTEGER,
-- 	WEIGHT_FOR_SET_4 INTEGER,
-- 	SET_5_REPS INTEGER,
-- 	WEIGHT_FOR_SET_5 INTEGER,
-- 	COMMENTS TEXT
-- )




BEGIN;

-- -- 1) Create table with constraints
-- CREATE TABLE IF NOT EXISTS EXERCISES (
--     exercise_name TEXT PRIMARY KEY,
--     type TEXT NOT NULL,
--     body_part TEXT NOT NULL,
--     -- Enforce allowed type values (lowercase)
--     CONSTRAINT exercises_type_chk CHECK (LOWER(type) IN ('cardio', 'strength', 'stretching'))
-- );

-- 2) Populate 40 rows
-- INSERT INTO EXERCISES (exercise_name, type, body_part) VALUES
--     -- Cardio (10)
--     ('Running',                  'cardio',   'legs'),
--     ('Jump Rope',                'cardio',   'full body'),
--     ('Cycling',                  'cardio',   'legs'),
--     ('Rowing (Erg)',             'cardio',   'back'),
--     ('Swimming',                 'cardio',   'full body'),
--     ('High Knees',               'cardio',   'legs'),
--     ('Burpees',                  'cardio',   'full body'),
--     ('Mountain Climbers',        'cardio',   'core'),
--     ('Stair Climbing',           'cardio',   'legs'),
--     ('Shadowboxing',             'cardio',   'upper body'),

    -- Strength (18)
    -- ('Back Squat',               'strength', 'legs'),
    -- ('Deadlift',                 'strength', 'posterior chain'),
    -- ('Bench Press',              'strength', 'chest'),
    -- ('Overhead Press',           'strength', 'shoulders'),
    -- ('Pull-Up',                  'strength', 'back'),
    -- ('Bent-Over Row',            'strength', 'back'),
    -- ('Hip Thrust',               'strength', 'glutes'),
    -- ('Forward Lunge',            'strength', 'legs'),
    -- ('Leg Press',                'strength', 'legs'),
    -- ('Standing Calf Raise',      'strength', 'calves'),
    -- ('Barbell Bicep Curl',       'strength', 'biceps'),
    -- ('Parallel Bar Dip',         'strength', 'triceps'),
    -- ('Plank Hold',               'strength', 'core'),
    -- ('Russian Twist',            'strength', 'core'),
    -- ('Kettlebell Swing',         'strength', 'posterior chain'),
    -- ('Farmer''s Carry',          'strength', 'full body'),
    -- ('Push-Up',                  'strength', 'chest'),
    -- ('Dumbbell Chest Fly',       'strength', 'chest'),

    -- -- Stretching (12)
    -- ('Standing Hamstring Stretch','stretching','hamstrings'),
    -- ('Quad Stretch',              'stretching','quads'),
    -- ('Calf Stretch',              'stretching','calves'),
    -- ('Hip Flexor Stretch',        'stretching','hip flexors'),
    -- ('Piriformis Stretch',        'stretching','glutes'),
    -- ('Child''s Pose',             'stretching','back'),
    -- ('Cat-Cow',                   'stretching','spine'),
    -- ('Cross-Body Shoulder Stretch','stretching','shoulders'),
    -- ('Overhead Triceps Stretch',  'stretching','triceps'),
    -- ('Doorway Chest Stretch',     'stretching','chest'),
    -- ('IT Band Stretch',           'stretching','legs'),
    -- ('Ankle Mobility Circles',    'stretching','ankles');

-- COMMIT;

-- SELECT * FROM EXERCISE_TRACKER;
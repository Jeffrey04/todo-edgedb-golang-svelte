CREATE MIGRATION m1kdgja4pidivk6cbraal344ujixyhxthprbtt5wywnguprs4qmeqa
    ONTO initial
{
  CREATE FUTURE nonrecursive_access_policies;
  CREATE TYPE default::Todo {
      CREATE REQUIRED PROPERTY created_at -> std::datetime;
      CREATE PROPERTY done_at -> std::datetime;
      CREATE REQUIRED PROPERTY item -> std::str;
      CREATE REQUIRED PROPERTY updated_at -> std::datetime;
  };
};

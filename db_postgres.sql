--
-- Table structure for table "follow"
--

CREATE TABLE "follow" (
  "id" SERIAL PRIMARY KEY,
  "userId" BIGINT NOT NULL,
  "userName" varchar(100) NOT NULL,
  "status" text,
  "followDate" TIMESTAMP WITH TIME ZONE NOT NULL,
  "unfollowDate" TIMESTAMP WITH TIME ZONE DEFAULT NULL,
  "lastAction" TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX ON "follow" ("userId");

--
-- Table structure for table "tweet"
--

CREATE TABLE "tweet" (
  "id" SERIAL PRIMARY KEY,
  "content" text NOT NULL,
  "date" TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX ON "tweet" ("content");

--
-- Table structure for table "reply"
--

CREATE TABLE "reply" (
  "id" SERIAL PRIMARY KEY,
  "userId" BIGINT NOT NULL,
  "userName" VARCHAR(100) NOT NULL,
  "tweetId" BIGINT NOT NULL,
  "status" TEXT NOT NULL,
  "answer" TEXT NOT NULL,
  "replyDate" TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX ON "reply" ("tweetId");

--
-- Table structure for table "favorite"
--

CREATE TABLE "favorite" (
  "id" SERIAL PRIMARY KEY,
  "userId" BIGINT NOT NULL,
  "userName" VARCHAR(100) NOT NULL,
  "tweetId" BIGINT NOT NULL,
  "status" TEXT NOT NULL,
  "favDate" TIMESTAMP WITH TIME ZONE NOT NULL,
  "unfavDate" TIMESTAMP WITH TIME ZONE NULL,
  "lastAction" TIMESTAMP WITH TIME ZONE NOT NULL
);

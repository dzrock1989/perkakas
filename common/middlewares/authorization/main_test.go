package authorization

import (
	"os"
	"testing"

	"github.com/tigapilarmandiri/perkakas/common/db"
	"github.com/tigapilarmandiri/perkakas/common/util"
	"github.com/tigapilarmandiri/perkakas/configs"
	"gorm.io/gorm"
)

var dbGorm *gorm.DB

func TestMain(m *testing.M) {
	confOpt := configs.ConfigOpts{
		EnvFile: "../../../.env",
	}
	configs.LoadConfigsWithOption(confOpt)

	dbConn := &db.DBConn{
		Info:       configs.Config.DB,
		SilentMode: true,
	}

	dbConn.Info.Name += "_test"

	var err error
	dbGorm, err = dbConn.Open()
	if err != nil {
		util.Log.Fatal().Msg(err.Error())
	}

	_, err = dbGorm.Raw(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`).Rows()
	if err != nil {
		util.Log.Error().Msg(err.Error())
	}

	err = dbGorm.Exec(`CREATE TABLE if not exists "wilayahs" (
  "id" uuid NOT NULL,
  "name" text COLLATE "pg_catalog"."default" NOT NULL,
	"jenis" text,
  "parent_id" uuid,
	"active" boolean,
	"deleted_at" timestamptz
)
;`).Error
	if err != nil {
		panic(err)
	}

	var count int
	dbGorm.Raw("select count(1) from wilayahs").Scan(&count)
	if count == 0 {
		err = dbGorm.Exec(`ALTER TABLE "public"."wilayahs" ADD CONSTRAINT "wilayah_pkey" PRIMARY KEY ("id");`).Error
		if err != nil {
			panic(err)
		}

		// err = dbGorm.Exec(`ALTER TABLE "public"."wilayahs" ADD CONSTRAINT "wilayah_parent_id_fkey" FOREIGN KEY ("parent_id") REFERENCES "public"."wilayahs" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;`).Error
		// if err != nil {
		// 	panic(err)
		// }

		err = dbGorm.Exec(`BEGIN;
	INSERT INTO "public"."wilayahs" ("id", "name", "parent_id", "active") VALUES ('22d72fbe-3dee-4044-884e-75f96cc1a27b', 'Indonesia', NULL, true);
	INSERT INTO "public"."wilayahs" ("id", "name", "parent_id", "active") VALUES ('2843b090-fa0c-488b-a6e4-a0d1c336dc0c', 'Jakarta', '22d72fbe-3dee-4044-884e-75f96cc1a27b', true);
	INSERT INTO "public"."wilayahs" ("id", "name", "parent_id", "active") VALUES ('c57dcfef-6fe0-4522-ac91-005fbe0021a1', 'Jakarta Barat', '2843b090-fa0c-488b-a6e4-a0d1c336dc0c', true);
	INSERT INTO "public"."wilayahs" ("id", "name", "parent_id", "active") VALUES ('0d51076f-1fdf-40eb-9bc0-1f939b75f376', 'Jakarta Selatan', '2843b090-fa0c-488b-a6e4-a0d1c336dc0c', true);
	INSERT INTO "public"."wilayahs" ("id", "name", "parent_id", "active") VALUES ('e64bd782-605b-49b9-b4b0-805fbc358ff6', 'Nusa Tenggara Barat', '22d72fbe-3dee-4044-884e-75f96cc1a27b', true);
	INSERT INTO "public"."wilayahs" ("id", "name", "parent_id", "active") VALUES ('63d9fefc-f566-4ff4-90f6-af42f57a6ab1', 'Lombok Tengah', 'e64bd782-605b-49b9-b4b0-805fbc358ff6', true);
	INSERT INTO "public"."wilayahs" ("id", "name", "parent_id", "active") VALUES ('a0311763-e68d-4621-85c4-73601d4a52d0', 'Praya', '63d9fefc-f566-4ff4-90f6-af42f57a6ab1', true);
	INSERT INTO "public"."wilayahs" ("id", "name", "parent_id", "active") VALUES ('0c0477e6-341d-4062-afad-8cd0f37f8b27', 'Leneng', '63d9fefc-f566-4ff4-90f6-af42f57a6ab1', true);
	INSERT INTO "public"."wilayahs" ("id", "name", "parent_id", "active") VALUES ('03273c09-0f0f-4e59-bbfe-5750ad9bdbce', 'Meteng', '63d9fefc-f566-4ff4-90f6-af42f57a6ab1', true);
COMMIT;`).Error
		if err != nil {
			panic(err)
		}

	}

	err = InitPreparedStatements(dbGorm)
	if err != nil {
		util.Log.Fatal().Msg(err.Error())
	}

	codeInt := m.Run()

	err = dbGorm.Exec("drop table wilayahs").Error
	if err != nil {
		panic(err)
	}

	os.Exit(codeInt)
}

# Generated by Django 4.2.16 on 2024-12-09 08:51

from django.db import migrations, models


class Migration(migrations.Migration):

    dependencies = [
        ("applications", "0013_applicationdeploymentmoduleorder"),
    ]

    operations = [
        migrations.AddField(
            model_name="application",
            name="app_tenant_id",
            field=models.CharField(
                help_text="应用对哪个租户的用户可用，当应用租户模式为全租户时，本字段值为空",
                max_length=32,
                null=True,
                verbose_name="应用租户 ID",
            ),
        ),
        migrations.AddField(
            model_name="application",
            name="app_tenant_mode",
            field=models.CharField(
                help_text="应用在租户层面的可用范围，可选值：全租户、指定租户",
                max_length=16,
                null=True,
                verbose_name="应用租户模式",
            ),
        ),
        migrations.AddField(
            model_name="application",
            name="tenant_id",
            field=models.CharField(
                db_index=True,
                help_text="本条数据的所属租户",
                max_length=32,
                null=True,
                verbose_name="租户 ID",
            ),
        ),
    ]
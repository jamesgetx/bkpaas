# Generated by Django 3.2.12 on 2024-04-09 07:37

from django.db import migrations, models
import django.db.models.deletion
import paasng.platform.mgrlegacy.entities
import paasng.platform.mgrlegacy.models
import paasng.utils.models


class Migration(migrations.Migration):

    dependencies = [
        ('applications', '0011_alter_application_type'),
        ('mgrlegacy', '0002_auto_20211130_1129'),
    ]

    operations = [
        migrations.CreateModel(
            name='CNativeMigrationProcess',
            fields=[
                ('id', models.AutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('region', models.CharField(help_text='部署区域', max_length=32)),
                ('created', models.DateTimeField(auto_now_add=True)),
                ('updated', models.DateTimeField(auto_now=True)),
                ('owner', paasng.utils.models.BkUserField(blank=True, db_index=True, max_length=64, null=True)),
                ('status', models.CharField(choices=[('default', 'Default'), ('on_migration', 'On migration'), ('migration_succeeded', 'Migration succeeded'), ('migration_failed', 'Migration failed'), ('confirmed', 'Confirmed'), ('on_rollback', 'On rollback'), ('rollback_succeeded', 'Rollback succeeded'), ('rollback_failed', 'Rollback failed')], default='default', max_length=20)),
                ('legacy_data', paasng.platform.mgrlegacy.models.DefaultAppLegacyDataField(default=paasng.platform.mgrlegacy.entities.DefaultAppLegacyData, help_text='记录迁移前的数据, 用于回滚')),
                ('created_at', models.DateTimeField(auto_now_add=True, help_text='操作记录的创建时间')),
                ('confirm_at', models.DateTimeField(help_text='用户确认的时间', null=True)),
                ('details', paasng.platform.mgrlegacy.models.ProcessDetailsField(default=paasng.platform.mgrlegacy.entities.ProcessDetails, help_text='迁移过程详情')),
                ('app', models.ForeignKey(db_constraint=False, on_delete=django.db.models.deletion.CASCADE, to='applications.application')),
            ],
            options={
                'ordering': ['-created_at'],
                'get_latest_by': 'created_at',
            },
        ),
    ]

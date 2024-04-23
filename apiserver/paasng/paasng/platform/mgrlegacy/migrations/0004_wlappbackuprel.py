# Generated by Django 3.2.12 on 2024-04-23 07:31

from django.db import migrations, models
import django.db.models.deletion


class Migration(migrations.Migration):

    dependencies = [
        ('applications', '0011_alter_application_type'),
        ('mgrlegacy', '0003_cnativemigrationprocess'),
    ]

    operations = [
        migrations.CreateModel(
            name='WlAppBackupRel',
            fields=[
                ('id', models.AutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('region', models.CharField(help_text='部署区域', max_length=32)),
                ('created', models.DateTimeField(auto_now_add=True)),
                ('updated', models.DateTimeField(auto_now=True)),
                ('original_id', models.UUIDField(verbose_name='原 WlApp uuid')),
                ('backup_id', models.UUIDField(verbose_name='对应备份的 WlApp uuid')),
                ('app_environment', models.OneToOneField(db_constraint=False, on_delete=django.db.models.deletion.CASCADE, to='applications.applicationenvironment')),
            ],
            options={
                'abstract': False,
            },
        ),
    ]

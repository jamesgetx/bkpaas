# Generated by Django 3.2.25 on 2024-10-17 13:04

from django.db import migrations, models
import django.db.models.deletion


class Migration(migrations.Migration):

    dependencies = [
        ('applications', '0012_application_is_ai_agent_app'),
        ('bkapp_model', '0015_auto_20240913_1509'),
    ]

    operations = [
        migrations.CreateModel(
            name='ProcessServicesFlag',
            fields=[
                ('id', models.AutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('region', models.CharField(help_text='部署区域', max_length=32)),
                ('created', models.DateTimeField(auto_now_add=True)),
                ('updated', models.DateTimeField(auto_now=True)),
                ('auto_created', models.BooleanField(default=False, verbose_name='是否自动创建')),
                ('app_environment', models.OneToOneField(db_constraint=False, on_delete=django.db.models.deletion.CASCADE, to='applications.applicationenvironment')),
            ],
            options={
                'abstract': False,
            },
        ),
    ]
// export TF_VAR_access_key=Your_access_key
// export TF_VAR_secret_key=Your_access_key
variable "access_key" {}
variable "secret_key" {}
variable "region" {
  default = "ap-northeast-1"
}
variable "cluster_name" {}

provider "aws" {
  profile = "default"
  access_key = var.access_key
  secret_key = var.secret_key
  region = var.region
}


// CTF kosenctfx用のECSを設定する
// 1. VPCを立ち上げる

resource "aws_vpc" "kosenctfx" {
  cidr_block = "10.0.0.0/16"
  tags = {
    Name = "Challenges VPC"
  }
}

// VPCは1つ以上のSubnetを持つ
// ここではPublic Subnetを2つ定義している
//  Private Subnetだとデフォルトでインターネットに接続不可能になっていて
//  接続するためにはNAT Gatewayを設定することになる
//  今回は基本的に外部接続可能ということにしてPublic Subnetのみを使用する

resource "aws_subnet" "a" {
  vpc_id = aws_vpc.kosenctfx.id
  cidr_block = "10.0.1.0/24"
  availability_zone = "${var.region}a"
  tags = {
    Name = "Challenge VPC Subnet A"
  }
}

resource "aws_subnet" "c" {
  vpc_id = aws_vpc.kosenctfx.id
  cidr_block = "10.0.2.0/24"
  availability_zone = "${var.region}a"
  tags = {
    Name = "Challenge VPC Subnet C"
  }
}

// Public Subnetが外部に接続するためのInternet GatewayとRoute Tableを作成する
// 作成してVPCに紐付けるだけ。デフォルト設定からいじるところはない
// Route TableはあとでSubnetに紐付ける

resource "aws_internet_gateway" "kosenctfx" {
  vpc_id = aws_vpc.kosenctfx.id
  tags = {
    Name = "Challenge VPC Internet Gateway"
  }
}

resource "aws_route_table" "kosenctfx" {
  vpc_id = aws_vpc.kosenctfx.id

  // (Subnet以外への)全てのトラフィックは外部に向く
  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.kosenctfx.id
  }

  tags = {
    Name = "Challenge VPC Route Table"
  }
}

// Route TableとSubnetを紐付ける
// これで10.0.1.0/24と10.0.2.0/24がそれぞれのSubnetに向くはず

resource "aws_route_table_association" "a" {
  subnet_id = aws_subnet.a.id
  route_table_id = aws_route_table.kosenctfx.id
}

resource "aws_route_table_association" "c" {
  subnet_id = aws_subnet.c.id
  route_table_id = aws_route_table.kosenctfx.id
}

// 2. ECSで使用するIAM RoleやPolicyを作成する


// 3. ECS Clusterを作成する
// ECS Clusterと、Clusterに属するInstanceのAuto Scalingルールを決定するCapacity Providerを作成する
// 正直良くわからないけどいい感じに動作するAuto Scalingがほしかったので、外部moduleを使わせてもらった
data "aws_ami" "ecs_ami" {
  most_recent = true
  owners      = ["amazon"]

  filter {
    name   = "name"
    values = ["amzn-ami-*-amazon-ecs-optimized"]
  }
}

module "app_ecs_cluster" {
  source = "github.com/trussworks/terraform-aws-ecs-cluster"

  name        = var.cluster_name
  environment = "dev"

  image_id      = data.aws_ami.ecs_ami.image_id
  instance_type = "t2.micro"

  vpc_id = aws_vpc.kosenctfx.id
  subnet_ids       = [ aws_subnet.a.id, aws_subnet.c.id ]
  max_size         = 10
  desired_capacity = 1
  min_size         = 0
}

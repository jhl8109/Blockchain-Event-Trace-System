# Blockchain-Event-Trace-System

### 한국컴퓨터종합학술대회 2023 (KCC 2023)

## 개요
허가된 블록체인은 공급망 관리, 의료, 투표 플랫폼 등 다양한 애플리케이션에 안전하고 투명한 솔루션을 제공할 수 있어 주목받고 있다. 

> #### 문제점) 그러나 블록체인 네트워크에 쌓이는 데이터의 양이 증가할수록 데이터에 대한 조회, 트랜잭션 수행 등의 성능 저하 문제가 발생한다.
> #### 해결방안) 이를 해결하기 위한 방안으로 본 논문에서는 이벤트 데이터를 통해 런타임 데이터를 수집한다. 또한, 데이터를 시간적, 공간적, 참여자, 트랜잭션 데이터로 분류하여 오프체인 데이터를 형성하는 방안을 제시한다. 

블록체인 실행환경에서 발생하는 이벤트를 수신 및 처리하는 메커니즘과 오프체인 데이터 분류를 제시함으로써 허가형 블록체인에서 오프체인 데이터를 구축하기 위한 기반 템플릿 모델로 본 방법이 적용될 수 있다. 
<br> 사례연구로 Hyperledger Fabric에 구축된 P2P 중고품 거래 플랫폼을 제시하여 제안된 시스템이 활용 가능함을 보이고 오프체인 데이터 조회할 경우 회당 약 29ms 정도의 성능이 개선됨을 확인하였다.

## 오프체인 데이터 분류 

허가형 블록체인은 액세스 제어, 높은 개인정보 보호, 다양한 산업에 대한 범용성을 갖는다. 이러한 특징을 반영하기 위해 아래 그림과 같이 데이터를 분류하였다.
<br> 제시한 데이터 분류는 데이터 분류별 모니터링과 분석을 통해 시스템(또는 특정 도메인)에 대한 추세와 패턴을 식별할 수 있다. 또한, 필요에 따라 데이터 분류들을 Join하여 새로운 인사이트를 얻는 데 활용할 수 있다.

<div align=center> <img width="500" src="https://github.com/jhl8109/Restful-Service/assets/78259314/5e9be185-7cca-41cc-9646-c7762b80dfce"/> </div>

## 아키텍처


<div align=center> <img width="500" src="https://github.com/jhl8109/Restful-Service/assets/78259314/50f44a7c-aa3d-4adc-9777-de186d8f4b80"/> </div>
 





## 사례연구
 <div align=center> 
  <img width="500" src="https://github.com/jhl8109/Restful-Service/assets/78259314/9d35fc6c-450f-40bd-8385-f2aad282a9d8"/>
  <img width="500" src="https://github.com/jhl8109/Restful-Service/assets/78259314/a0af3540-c8f3-4d87-a155-dbd46a737bf1"/>
</div>

## 성능평가
<div align=center>
  
![스크린샷 2023-05-23 오전 12 06 04](https://github.com/jhl8109/Restful-Service/assets/78259314/f7d711c3-0c62-421e-9b6b-184582ff45d0)
![스크린샷 2023-05-23 오전 12 05 51](https://github.com/jhl8109/Restful-Service/assets/78259314/88e90651-4b3a-4398-b317-5650ed61e6c8)
</div>

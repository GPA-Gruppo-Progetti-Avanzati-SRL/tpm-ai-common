       IDENTIFICATION DIVISION.
       PROGRAM-ID.    YPBCEPGP.
       AUTHOR.        TAS.
      *================================================================*
      *                                                                *
      *      NNNN    NNN         CCCCCCC        HHH   HHH              *
      *      NNNNN   NNN        CCCCCCCC        HHH   HHH              *
      *      NNNNNN  NNN        CCC             HHHHHHHHH              *
      *      NNN NNNNNNN        CCC             HHHHHHHHH              *
      *      NNN   NNNNN  ...   CCCCCCCC  ...   HHH   HHH  ...         *
      *      NNN     NNN  ...    CCCCCCC  ...   HHH   HHH  ....        *
      *                                                                *
      *================================================================*
      *      NETWORK           COMPUTER        HOUSE      - BOLOGNA -  *
      *================================================================*
      *  VERSIONE 01.00 DEL : 07/05/15 --- ULTIMA MODIFICA : XX/XX/XX  *
      *  SCHEDA ........... : R05415                                   *
      *================================================================*
      *                                                                *
      * YPBCEPGP:  GESTISCE IL FLUSSO PGPF DEDICATO ALLA CONTABILITA'  *
      *            ALLA CONTABILITA' VERSO I CONTI CORRENTI DEGLI      *
      *            ESERCENTI, ESEGUENDO:                               *
      *                                                                *
      *            1)LA CREAZIONE FLUSSO CON TRACCIATO CRVFSD50        *
      *              DEDICATO ALLA CONTABILITA' VERSO I                *
      *              CONTI CORRENTI DEGLI ESERCENTI.                   *
      *                                                                *
      *            2)IL CARICAMENTO TABELLA YPTBPGPF CON MESSAGGI      *
      *              DEDICATI ALLA CONTABILITA' VERSO I CONTI          *
      *              CORRENTI DEGLI ESERCENTI (LL = 2000).             *
      *              VENGONO ELABORATI, AL FINE DEL CARICAMENTO,       *
      *              I MESSAGGI CHE HANNO:                             *
      *                                                                *
      *              . MSG-TYPE-ID = 0610                              *
      *              . FUNCT-CODE  = TUTTI                             *
      *                                                                *
      *                                                                *
      *================================================================*
       ENVIRONMENT DIVISION.
      *================================================================*
       CONFIGURATION SECTION.
       SPECIAL-NAMES.
           DECIMAL-POINT IS COMMA.
      *
       INPUT-OUTPUT SECTION.
       FILE-CONTROL.
           SELECT IPAYMENT               ASSIGN       TO IPAYMENT
                                         FILE STATUS  IS ST-IPAYMENT.

           SELECT OPECONT                ASSIGN       TO OPECONT
                                         FILE STATUS  IS ST-OPECONT.

FIANNH     SELECT OPECONTB                ASSIGN       TO OPECONTB
FIANNH                                   FILE STATUS  IS ST-OPECONTB.

R14316     SELECT XYDCONT                ASSIGN       TO XYDCONT
R14316                                   FILE STATUS  IS ST-XYDCONT.
R14316
R14316     SELECT OUTDCD                 ASSIGN       TO OUTDCD
R14316                                   FILE STATUS  IS ST-OUTDCD.
R14316
           SELECT OUSCARTI               ASSIGN       TO OUSCARTI
                                         FILE STATUS  IS ST-OUSCARTI.

R11817     SELECT OUSCART2               ASSIGN       TO OUSCART2
R11817                                   FILE STATUS  IS ST-OUSCART2.

R05818     SELECT OUTFXML                ASSIGN       TO OUTFXML
R05818                                   FILE STATUS  IS ST-OUTFXML.

R11422     SELECT OUTFXM2                ASSIGN       TO OUTFXM2
R11422                                   FILE STATUS  IS ST-OUTFXM2.

R12019     SELECT BILLCCB                ASSIGN       TO BILLCCB
R12019                                   FILE STATUS  IS ST-BILLCCB.

           SELECT YYDTABE                ASSIGN       TO YYDTABE
                                         ORGANIZATION IS INDEXED
                                         RECORD KEY   IS YYDTABE-KEY
R12117*                                  ACCESS MODE  IS RANDOM
R12117                                   ACCESS MODE  IS DYNAMIC
                                         FILE STATUS  IS ST-YYDTABE.

           SELECT ST                     ASSIGN       TO STAMPA.

           SELECT YPOERRO                ASSIGN       TO YPOERRO
                                         FILE STATUS  IS ST-YPOERRO.

      *================================================================*
       DATA DIVISION.
      *
       FILE SECTION.
      *
      *----------------------------------------------------------------*
      * IPAYMENT      :  FLUSSO DI IPAYMENT                    (INPUT) *
      *----------------------------------------------------------------*
       FD  IPAYMENT                      LABEL RECORD STANDARD
                                         RECORDING MODE IS F
                                         BLOCK      0 RECORDS.

       01  IPAYMENT-REC                  PIC  X(2000).
      *
      *----------------------------------------------------------------*
      * OPECONT       :  OPERAZIONI CONTABILI                 (OUTPUT) *
      *----------------------------------------------------------------*
       FD  OPECONT                       LABEL RECORD STANDARD
                                         RECORDING MODE IS F
                                         BLOCK      0 RECORDS.

       01  OPECONT-REC.
           03  FILLER                    PIC X(478).

FIANNH*----------------------------------------------------------------*
FIANNH* OPECONT       :  OPERAZIONI CONTABILI                 (OUTPUT) *
FIANNH*----------------------------------------------------------------*
FIANNH FD  OPECONTB                      LABEL RECORD STANDARD
FIANNH                                   RECORDING MODE IS F
FIANNH                                   BLOCK      0 RECORDS.
FIANNH
FIANNH 01  OPECONTB-REC.
FIANNH     03  FILLER                    PIC X(478).

R14316*----------------------------------------------------------------*
R14316* XYDCONT      :  SEQUENZIALE PER ADDEBITI CONTABILI NORMALIZZATI*
R14316*----------------------------------------------------------------*
R14316 FD  XYDCONT                       LABEL RECORD STANDARD
R14316                                   BLOCK    0   RECORDS.
R14316
R14316 01  XYDCONT-REC                   PIC  X(500).
R14316
R14316*----------------------------------------------------------------*
R14316* OUTDCD       :  ARCHIVIO DI OUTPUT CONTENENTE I MOVIMENTI DCD  *
R14316*----------------------------------------------------------------*
R14316 FD  OUTDCD                        LABEL RECORD STANDARD
R14316                                   BLOCK    0   RECORDS.
R14316
R14316 01  REC-OUTDCD                    PIC  X(800).
R14316
      *----------------------------------------------------------------*
      * OUSCARTI      :  OPERAZIONI CONTABILI SCARTATE        (OUTPUT) *
      *----------------------------------------------------------------*
       FD  OUSCARTI                      LABEL RECORD STANDARD
                                         RECORDING MODE IS F
                                         BLOCK      0 RECORDS.

       01  OUSCARTI-REC.
           03  FILLER                    PIC X(2000).

R11817*----------------------------------------------------------------*
R11817* OUSCART2 :  OPERAZIONI CONTABILI SCARTATE  Postepay Evo Business
R11817*----------------------------------------------------------------*
R11817 FD  OUSCART2                      LABEL RECORD STANDARD
R11817                                   RECORDING MODE IS F
R11817                                   BLOCK      0 RECORDS.
R11817
R11817 01  OUSCART2-REC.
R11817     03  OUSCART2-PARTE1           PIC X(2000).
R11817     03  OUSCART2-PARTE2           PIC X(0800).
R11817
      *----------------------------------------------------------------*
R05818* OUTFXML       :  DATI PER PRODUZIONE FLUSSO XML                *
R05818*----------------------------------------------------------------*
R05818 FD  OUTFXML                       LABEL RECORD STANDARD
R05818                                   RECORDING MODE IS F
R05818                                   BLOCK      0 RECORDS.
R05818
R05818 01  OUTFXML-REC.
R05818     03  OUTFXML-TIPO-REC          PIC X(01).
R05818     03  OUTFXML-DATI              PIC X(192).
R05818     03  OUTFXML-DESC              PIC X(140).
R05818     03  OUTFXML-FILLER            PIC X(17).
R05818
      *----------------------------------------------------------------*
R11422* OUTFXM2       :  DATI PER PRODUZIONE FLUSSO XML                *
R11422*----------------------------------------------------------------*
R11422 FD  OUTFXM2                       LABEL RECORD STANDARD
R11422                                   RECORDING MODE IS F
R11422                                   BLOCK      0 RECORDS.
R11422
R11422 01  OUTFXM2-REC                   PIC X(500).
      *----------------------------------------------------------------*
      * YYDTABE       :  TABELLA VSAM                         (INPUT)  *
      *----------------------------------------------------------------*
       FD  YYDTABE.

       01  YYDTABE-REC-FD.
           03 YYDTABE-KEY.
              05 YYDTABE-COD             PIC X(003).
              05 YYDTABE-COD-VAR         PIC X(003).
              05 YYDTABE-PROG            PIC 9(003).
              05 YYDTABE-TIPO-CRT        PIC X(002).
              05 YYDTABE-CAUSALE         PIC X(010).
              05 YYDTABE-RESTO-KEY       PIC X(009).
           03 YYDTABE-DATI               PIC X(500).
           03 FILLER                     PIC X(1470).

      *----------------------------------------------------------------*
      * ST            :  FILE STAMPA                          (OUTPUT) *
      *----------------------------------------------------------------*
       FD  ST                            LABEL RECORD STANDARD
                                         RECORDING MODE F
                                         BLOCK 0      RECORDS.
       01  REC-ST.
           03 FILLER                     PIC X(01).
           03 RIGA-ST                    PIC X(132).

      *================================================================*
      *    Sequenziale di output record errati per stampa              *
      *================================================================*
       FD  YPOERRO                   RECORDING MODE IS F
                                     LABEL RECORD STANDARD
                                     BLOCK    0   RECORDS.

       01  YPOERRO-REC                    PIC  X(0080).
      *
      *================================================================*
R12019* BILLCCB       :  FILE DI OUT FLUSSO BILL CCB  (SEQUENZIALE)    *
R12019*================================================================*
R12019 FD  BILLCCB                   RECORDING MODE IS F
R12019                               LABEL RECORD STANDARD
R12019                               BLOCK    0   RECORDS.
R12019
R12019 01  BILLCCB-REC               PIC  X(0300).
R12019*

      *================================================================*
       WORKING-STORAGE SECTION.
       01  PROGRAMMA                     PIC  X(08)  VALUE 'YPBCEPGP'.
R11422 01  W100-PGM-CALL                 PIC  X(08) VALUE SPACES.
       01  CRVYD228                      PIC  X(08)  VALUE 'CRVYD228'.
       01  XYRC0005                      PIC  X(08)  VALUE 'XYRC0005'.
       01  YYUCDATA                      PIC  X(08)  VALUE 'YYUCDATA'.
       01  Z3BCGE90                      PIC  X(08)  VALUE 'Z3BCGE90'.
R14316 01  Z3BCUIFA                      PIC  X(08)  VALUE 'Z3BCUIFA'.
R15420 01  Z3BCUI99                      PIC  X(08)  VALUE 'Z3BCUI99'.
R05818 01  WK-ACS108BT                    PIC X(008) VALUE 'ACS108BT'.
      *
       01  ST-IPAYMENT                   PIC  X(02)  VALUE '00'.
           88  IPAY-NORMAL               VALUE '00'.
           88  IPAY-EOF                  VALUE '10'.
       01  ST-OPECONT                    PIC  X(02).
           88  OPEC-NORMAL               VALUE '00'.
FIANNH 01  ST-OPECONTB                   PIC  X(02).
FIANNH     88  OPECB-NORMAL              VALUE '00'.
R14316 01  ST-XYDCONT                    PIC  X(02).
R14316     88  XYDCONT-NORMAL            VALUE '00'.
R14316 01  ST-OUTDCD                     PIC  X(02).
R14316     88  OUTDCD-NORMAL             VALUE '00'.
       01  ST-OUSCARTI                   PIC  X(02).
           88  OUSC-NORMAL               VALUE '00'.
R12019 01  ST-BILLCCB                    PIC  X(02).
R12019     88  BILC-NORMAL               VALUE '00'.
R05818 01  ST-OUTFXML                    PIC  X(02).
R05818     88  OXML-NORMAL               VALUE '00'.
R11422 01  ST-OUTFXM2                    PIC  X(02).
R11422     88  OXM2-NORMAL               VALUE '00'.
R11817 01  ST-OUSCART2                   PIC  X(02).
R11817     88  OUS2-NORMAL               VALUE '00'.
       01  ST-YPOERRO                    PIC  X(02).
           88  OERR-NORMAL               VALUE '00'.
       01  ST-YYDTABE                    PIC  X(02).
           88  YYDT-NORMAL               VALUE '00'.
           88  YYDT-NOTFND               VALUE '23'.
R12117 01  WK-TABE-SMAC                  PIC  X(02).
R12117     88  FINE-TABE-SMAC            VALUE 'SI'.
      *-----------------------------------------------
TK1274 01  WK-TROVATO-SU-FPR             PIC  X(02).
TK1274     88  TROVATO-SU-FPR-SI         VALUE 'SI'.
TK1274     88  TROVATO-SU-FPR-NO         VALUE 'NO' '  '.
      *----------  CAMPI DI COMODO PER DATA E ORA
R11422 01  W100-DATA-SOLARE9                 PIC 9(08) VALUE ZEROES.
R11422 01  W100-DATA-SOLARE.
R11422   03  W100-DATA-SOLARE-SS             PIC X(02) VALUE '20'.
R11422   03  W100-DATA-SOLARE-AAMMGG         PIC X(06) VALUE SPACES.
      *
      *================================================================*
      *    Area con valori fissi                                       *
      *================================================================*
       01  WS-AREA-VALO-FISSI.
      *-
           03  FILLER                       PIC X.
      *
           03  W300-SCART-DUPKEY            PIC 9(11).
R14316*
R14316     03  W100-ROUT-DECO               PIC X(08) VALUE 'YPRCP008'.
      *
           03  WS-MCC-PEDAGGI               PIC  9(04) VALUE 4784.
           03  WS-MCC-CSD                   PIC  9(04) VALUE 6010.
           03  WS-MCC-ATM                   PIC  9(04) VALUE 6011.
      *
           03  WS-SCRI-YPOE                 PIC 9(01)  VALUE 0.
               88 WS-SCRI-YPOE-NO           VALUE 0.
               88 WS-SCRI-YPOE-SI           VALUE 1.
      *
R11422     03  WS-LETTO-FAD                 PIC 9(01)  VALUE 0.
R11422         88 WS-LETTO-FAD-NO           VALUE 0.
R11422         88 WS-LETTO-FAD-SI           VALUE 1.
      *
           03  VISTO-REC-TESTA-681          PIC  9(01) VALUE 0.
               88 VISTO-REC-TESTA-681-NO    VALUE 0.
               88 VISTO-REC-TESTA-681-SI    VALUE 1.
      *
           03  WS-PRIMA-VOLTA               PIC 9(01) VALUE 0.
               88 WS-PRIMA-VOLTA-SI         VALUE 0.
               88 WS-PRIMA-VOLTA-NO         VALUE 1.
      *
           03  VALORE-MAX                   PIC S9(4) COMP VALUE +20.
           03  YP-INDMAX                    PIC S9(4) COMP VALUE +5.
      *
           03  COM-ACTION-CODE              PIC 9(03) VALUE 0.
               88  TESTA-680                VALUE 680.
               88  TESTA-681                VALUE 681.
      *
           03  RIGA-DI-CUI.
               05  FILLER-YY                PIC X(08) VALUE 'DI CUI .'.
               05  RIGA-DI-CUI-RESTO        PIC X(28).
               05  RIGA-DI-CUI-NUM          PIC ZZZ.ZZZ.ZZ9.
               05  FILLER                 PIC X(11) VALUE '  IMPORTO: '.
               05  RIGA-DI-CUI-IMP          PIC ZZZ.ZZZ.ZZ9,99.
      *
           03  WS-AREA-TEST-YPOE.
               05  WS-AREA-TEST-YPOE-DESC       PIC X(0035)
                   VALUE '     Descrizione dell''errore     '.
R14316*        05  WS-AREA-TEST-YPOE-PAN        PIC X(0020)
R14316*            VALUE '        Payment uid '.
R11422*        05  WS-AREA-TEST-YPOE-IBAN       PIC X(0035)
R11422*            VALUE '        Iban                      '.
R11422         05  WS-AREA-TEST-YPOE-NRFA       PIC X(0015)
R11422             VALUE '  N.Rapporto  '.
R11422         05  WS-AREA-TEST-YPOE-NRFA       PIC X(0015)
R11422             VALUE 'Merchant ID   '.
R14316*        05  WS-AREA-TEST-YPOE-TYPE       PIC X(0025)
R14316*            VALUE '    Type Account         '.
R14316         05  WS-AREA-TEST-YPOE-IMPO       PIC X(0010)
R14316             VALUE 'Importo   '.
      *
      *================================================================*
      *    Area inizializzata a low-value                              *
      *================================================================*
       01  WS-LOW-VALUE.
           03  FILLER                       PIC X.
      *
      *================================================================*
      *    Area inizializzata a zero                                   *
      *================================================================*
       01  WS-ZERO.
      *-
           03  FILLER                       PIC 9(01).
R11422     03  WS-APPO-IMPO                 PIC  9(12).
R11422     03  WS-APPO-IMPO-R REDEFINES WS-APPO-IMPO
R11422                                      PIC 9(10)V99.
R04417*
R04417     03  WK-DATA-NUM                  PIC  9(08).
      *
R05818     03  WS-PAYMT-GEN-DATE            PIC 9(08).
R05818     03  WS-BUSINESS-DATE             PIC 9(08).
R05818     03  WS-VALUE-DATE                PIC 9(08).
R05818     03  WS-REFER-PERIOD-DATE-200     PIC 9(08).
      *
           03  WK-END                       PIC  9(02).
R14316     03  WS-MONTE-MONETA-V            PIC 9(09)V9(02).
R14316     03  WS-SALDO-D                   PIC +Z(012)9,99.
R14316     03  WS-IMPORTO-TOT               PIC S9(13)V99.
R14316     03  WS-IMPORTO-MOV               PIC S9(10)V99.
R14316     03  WS-IMPORTO-MOV-D             PIC +Z(10),ZZ.
R05818     03  WS-TOTALI-IMPO-CC            PIC 9(15)V9(3).
R11422     03  WS-TOTALI-IMPO-CC-B          PIC 9(15)V9(3).
R05818     03  WS-NUM-OPER-CC               PIC 9(09).
R11422     03  WS-NUM-OPER-CC-B             PIC 9(09).
      *
R14316     03  W100-PROGR                   PIC 9(10).
      *
           03  COMSD50-RAPPORT              PIC S9(12) COMP-3.
      *
R14316     03  WS-APPO-ALIAS                PIC ZZZZZZZZ9.
      *
           03  ULTI-RECO                    PIC 9(01).
               88 ULTI-RECO-TEST-680        VALUE 1.
               88 ULTI-RECO-TEST-681        VALUE 2.
               88 ULTI-RECO-CODA            VALUE 3.
               88 ULTI-RECO-DETT            VALUE 4.
R03817*
R03817     03  WS-IMPO-ZERO                 PIC 9(01).
R03817         88 WS-IMPO-ZERO-SI VALUE 0.
R03817         88 WS-IMPO-ZERO-NO VALUE 1.
      *
           03  WS-SALT-ELAB                 PIC 9(01).
               88 WS-SALT-ELAB-SI VALUE 0.
               88 WS-SALT-ELAB-NO VALUE 1.
      *
R05316     03  WS-SALT-CONT                 PIC 9(01).
R05316         88 WS-SALT-CONT-SI VALUE 0.
R05316         88 WS-SALT-CONT-NO VALUE 1.
      *
R00317     03  WS-TIPO-BILL                 PIC 9(01).
R00317         88 WS-TIPO-BILL-SI VALUE 0.
R00317         88 WS-TIPO-BILL-NO VALUE 1.
R11422*
R11422     03  WS-ELAB-CC-BANCARI           PIC 9(01).
R11422         88 WS-ELAB-CC-BANCARI-NO     VALUE 0.
R11422         88 WS-ELAB-CC-BANCARI-SI     VALUE 1.
      *-
           03  WS-TIPO-FUN-CODE             PIC 9(03).
               88  WS-TIPO-FUN-CODE-697     VALUE  697.
               88  WS-TIPO-FUN-CODE-695     VALUE  695.
               88  WS-TIPO-FUN-CODE-300     VALUE  300.
               88  WS-TIPO-FUN-CODE-VALI    VALUE  697
                                                   695
                                                   200
                                                   300
                                                  .
      *-
           03  WS-EOF-IPAYMENT              PIC 9(01).
           03  CTR-PROGRES                  PIC 9(11).
           03  CTR-CONT-LETTI-TOT           PIC 9(09).
           03  CTR-CONT-LETTI-SCAR          PIC 9(09).
R05316*    03  CTR-CONT-LETTI-SCAR-APO      PIC 9(09).
R05316*    03  CTR-CONT-LETTI-SCAR-GPO      PIC 9(09).
R05316*    03  CTR-CONT-LETTI-SCAR-MPO      PIC 9(09).
R05316*    03  CTR-CONT-LETTI-SCAR-POS      PIC 9(09).
R05316*    03  CTR-CONT-LETTI-SCAR-VPO      PIC 9(09).
R05316*    03  CTR-CONT-LETTI-SCAR-APOI     PIC 9(11).
R05316*    03  CTR-CONT-LETTI-SCAR-GPOI     PIC 9(11).
R05316*    03  CTR-CONT-LETTI-SCAR-MPOI     PIC 9(11).
R05316*    03  CTR-CONT-LETTI-SCAR-POSI     PIC 9(11).
R05316*    03  CTR-CONT-LETTI-SCAR-VPOI     PIC 9(11).
R05316     03  CTR-CONT-LETTI-SCAR-BILL     PIC 9(09).
R05316     03  CTR-CONT-LETTI-SCAR-BILLI    PIC 9(11).
           03  CTR-CONT-LETTI-HEAD-680      PIC 9(09).
           03  CTR-CONT-LETTI-HEAD-681      PIC 9(09).
           03  CTR-CONT-LETTI-DATI-200      PIC 9(09).
           03  CTR-CONT-LETTI-DATI-300      PIC 9(09).
R08421     03  CTR-CONT-LETTI-DATI-301      PIC 9(09).
           03  CTR-CONT-LETTI-TRAIL         PIC 9(09).
           03  CTR-CONT-SCARTI              PIC 9(09).
R05818     03  CTR-CONT-FXML                PIC 9(09).
R11422     03  CTR-CONT-FXM2                PIC 9(09).
R11817     03  CTR-CONT-SCART2              PIC 9(09).
           03  CTR-CONT-SCRITTI             PIC 9(09).
FIANNH     03  CTR-CONT-SCRITTI-B           PIC 9(09).
R14316     03  CTR-CONT-SCRITTI-CONT        PIC 9(09).
R14316     03  CTR-CONT-SCRITTI-DCD         PIC 9(09).
           03  CTR-CONT-SCRITTI-YPOERRO     PIC 9(09).
R12019     03  CTR-CONT-SCRITTI-BILLCCB     PIC 9(09).
           03  CTR-COM-IMPORTO              PIC 9(11)V99.
           03  CTR-TABPGPF-INSE             PIC 9(09).
R11422     03  CTR-TABFAS2-INSE             PIC 9(09).
TK1274     03  CTR-TABFPR-NOT-FOUND         PIC 9(09).
TK1274     03  CTR-TABFPR-LETTE             PIC 9(09).

           03  ETR-CONT-LETTI-TOT           PIC ZZZ.ZZZ.ZZ9.
           03  ETR-CONT-LETTI-SCAR          PIC ZZZ.ZZZ.ZZ9.
R05316*    03  ETR-CONT-LETTI-SCAR-APO      PIC ZZZ.ZZZ.ZZ9.
R05316*    03  ETR-CONT-LETTI-SCAR-GPO      PIC ZZZ.ZZZ.ZZ9.
R05316*    03  ETR-CONT-LETTI-SCAR-MPO      PIC ZZZ.ZZZ.ZZ9.
R05316*    03  ETR-CONT-LETTI-SCAR-POS      PIC ZZZ.ZZZ.ZZ9.
R05316*    03  ETR-CONT-LETTI-SCAR-VPO      PIC ZZZ.ZZZ.ZZ9.
R05316*    03  ETR-CONT-LETTI-SCAR-APOI     PIC ZZZ.ZZZ.ZZ9,99.
R05316*    03  ETR-CONT-LETTI-SCAR-GPOI     PIC ZZZ.ZZZ.ZZ9,99.
R05316*    03  ETR-CONT-LETTI-SCAR-MPOI     PIC ZZZ.ZZZ.ZZ9,99.
R05316*    03  ETR-CONT-LETTI-SCAR-POSI     PIC ZZZ.ZZZ.ZZ9,99.
R05316*    03  ETR-CONT-LETTI-SCAR-VPOI     PIC ZZZ.ZZZ.ZZ9,99.
R05316     03  ETR-CONT-LETTI-SCAR-BILL     PIC ZZZ.ZZZ.ZZ9.
R05316     03  ETR-CONT-LETTI-SCAR-BILLI    PIC ZZZ.ZZZ.ZZ9,99.
           03  ETR-CONT-ELAB                PIC ZZZ.ZZZ.ZZ9.
           03  ETR-CONT-LETTI-HEAD-680      PIC ZZZ.ZZZ.ZZ9.
           03  ETR-CONT-LETTI-HEAD-681      PIC ZZZ.ZZZ.ZZ9.
           03  ETR-CONT-LETTI-DATI-200      PIC ZZZ.ZZZ.ZZ9.
           03  ETR-CONT-LETTI-DATI-300      PIC ZZZ.ZZZ.ZZ9.
R08421     03  ETR-CONT-LETTI-DATI-301      PIC ZZZ.ZZZ.ZZ9.
           03  ETR-CONT-LETTI-TRAIL         PIC ZZZ.ZZZ.ZZ9.
           03  ETR-CONT-SCARTI              PIC ZZZ.ZZZ.ZZ9.
R05818     03  ETR-CONT-FXML                PIC ZZZ.ZZZ.ZZ9.
R11422     03  ETR-CONT-FXM2                PIC ZZZ.ZZZ.ZZ9.
R11817     03  ETR-CONT-SCART2              PIC ZZZ.ZZZ.ZZ9.
           03  ETR-CONT-SCRITTI             PIC ZZZ.ZZZ.ZZ9.
FIANNH     03  ETR-CONT-SCRITTI-B           PIC ZZZ.ZZZ.ZZ9.
R14316     03  ETR-CONT-SCRITTI-CONT        PIC ZZZ.ZZZ.ZZ9.
R14316     03  ETR-CONT-SCRITTI-DCD         PIC ZZZ.ZZZ.ZZ9.
           03  ETR-CONT-SCRITTI-YPOERRO     PIC ZZZ.ZZZ.ZZ9.
R12019     03  ETR-CONT-SCRITTI-BILLCCB     PIC ZZZ.ZZZ.ZZ9.
           03  ETR-TABPGPF-INSE             PIC ZZZ.ZZZ.ZZ9.
R11422     03  ETR-TABFAS2-INSE             PIC ZZZ.ZZZ.ZZ9.
TK1274     03  ETR-TABFPR-LETTE             PIC ZZZ.ZZZ.ZZ9.
TK1274     03  ETR-TABFPR-NOT-FOUND         PIC ZZZ.ZZZ.ZZ9.
           03  IND1                         PIC S9(4).
           03  YP-IND                       PIC S9(4) COMP.
           03  ULT-TIPO-REC                      PIC  9(01).
               88 ULT-REC-TESTA-680              VALUE 0.
               88 ULT-REC-TESTA-681              VALUE 1.
               88 ULT-REC-CODA                   VALUE 2.
               88 ULT-REC-DATI                   VALUE 3.
           03  COM-DATE-TIME-H.
               05  COM-DATE-TIME-N          PIC 9(14).
      *
R05316     03  COM-DATE.
R05316         05  COM-DATE-N               PIC 9(08).
      *
           03  ESI-ERR-HEADER               PIC 9.
               88  HEADER-OK                VALUE 0.
               88  HEADER-ERR               VALUE 1.
           03  ESI-ERR-TABE                 PIC 9.
               88  TABE-OK                  VALUE 0.
               88  TABE-ERR                 VALUE 1.
           03  WS-SCRI-SCAR                 PIC 9.
               88  WS-SCRI-SCAR-NO          VALUE 0.
               88  WS-SCRI-SCAR-SI          VALUE 1.
R07420     03  WS-CONTAB                    PIC 9.
R07420         88  WS-CONTAB-NO             VALUE 0.
R07420         88  WS-CONTAB-SI             VALUE 1.
           03  SCRIVI-LAST-REC              PIC 9.
               88  SCRIVI-LAST-NO           VALUE 0.
               88  SCRIVI-LAST-SI           VALUE 1.
R14316     03  WS-IBAN-ATTI                 PIC 9.
R14316         88  WS-IBAN-ATTI-NO          VALUE 0.
R14316         88  WS-IBAN-ATTI-SI          VALUE 1.
R14316     03  WS-PAN-PPAY-EVOL-BUSI        PIC 9.
R14316         88  WS-PAN-PPAY-EVOL-BUSI-NO VALUE 0.
R14316         88  WS-PAN-PPAY-EVOL-BUSI-SI VALUE 1.
R05818     03  WS-TIPO-CC                   PIC 9.
R05818         88  WS-PP-EVO                VALUE 1.
R05818         88  WS-CC-POSTE              VALUE 2.
R05818         88  WS-CC-BANCARIO           VALUE 3.
R05818     03  WS-LETTA-GEP-CCB             PIC 9.
R05818         88  WS-LETTA-GEP-CCB-NO      VALUE 0.
R05818         88  WS-LETTA-GEP-CCB-SI      VALUE 1.
           03  WS-DATA-AAMMGG               PIC 9(6).
           03  WS-ORA-DAY                   PIC 9(08).
           03  WS-COM-RAPPORTO.
               05  WS-COM-RAPPORTON         PIC 9(12).
R12019     03  WK-IMPO-BILC                 PIC ZZZZZZZZZZZZ9,99.
R12019     03  WK-IMPO-BILC-NUME            PIC 9(10)V9(2).
TK1274     03  WK-IMPO-BILC-NUME-ASI        PIC 9(10)V9(2).
ZANCHI     03  WS-ORA.
               05  WS-ORA-HH                PIC X(02).
               05  WS-ORA-MM                PIC X(02).
               05  WS-ORA-SS                PIC X(02).
      *
      *================================================================*
      *    Area inizializzata a space                                  *
      *================================================================*
       01  WS-SPACE.
      *
           03  FILLER                       PIC X(01).
R05818     03  WS-DESC-MOV                  PIC X(140).
R14316     03  WS-MONTE-MONETA              PIC X(11).
R14316     03  WS-MONTE-MONETA-N REDEFINES WS-MONTE-MONETA  PIC  9(11).
R14316     03  WS-MONTE-MONETA-S            PIC X(01).
R14316*
R14316     03  W100-PROGR-APERTURA          PIC X(05).
R14316     03  W100-E-CAMPO.
R14316         05 FILLER                    PIC X(05).
R14316         05 W100-PROGR-ESA            PIC X(05).
R14316*
R14316     03 WS-PAN-EUC.
R14316       05  WS-BIN-EUC                 PIC X(06).
R14316       05  WS-CARTA-EUC               PIC X(09).
R14316       05  WS-CKD-EUC                 PIC X(01).
R14316       05  WS-FILLER-EUC              PIC X(03).
R14316*
R14316     03  WS-IBAN.
R14316       05  WS-IBAN-PAE                PIC X(02).
R14316       05  WS-IBAN-CTR                PIC X(02).
R14316       05  WS-RESTO-23.
R14316         07  WS-IBAN-CIN              PIC X(01).
R14316         07  WS-IBAN-ABI              PIC X(05).
R14316         07  WS-IBAN-CAB              PIC X(05).
R14316         07  WS-IBAN-CCC              PIC X(12).
R14316*
R05818     03 WS-DATI-XML-CODA.
R05818       05  WS-NUM-TRAN                PIC 9(08).
R05818       05  WS-TOTALE-IMP              PIC 9(15)V9(3).
R05818       05  WS-IBAN-MITT.
R05818          10  WS-IBAN-PAE-M              PIC X(02).
R05818          10  WS-IBAN-CTR-M              PIC X(02).
R05818          10  WS-RESTO-23-M.
R05818             15  WS-IBAN-CIN-M            PIC X(01).
R05818             15  WS-IBAN-ABI-M            PIC X(05).
R05818             15  WS-IBAN-CAB-M            PIC X(05).
R05818             15  WS-IBAN-CCC-M            PIC X(12).
R05818       05  FILLER                     PIC X(139).
R05818     03 WS-DATI-XML.
R05818       05  WS-IBAN-DEST               PIC X(27).
R05818       05  WS-RAGI-SOC                PIC X(60).
R05818       05  WS-IMPO-MOV                PIC 9(09)V99.
R05818       05  WS-INDIRIZZO               PIC X(35).
R05818       05  WS-CAP                     PIC X(5).
R05818       05  WS-LOC                     PIC X(30).
R05818       05  WS-PROV                    PIC X(2).
R05818       05  WS-NAZ                     PIC X(4).
R05818       05  WS-PAYMENT-UID             PIC 9(18).
      ******************************************************************
      *
R12117     03 WK-AREA-CAUS-OPE.
R12019*        05 WK-OLI-ELEMENTO-NW    OCCURS 200.
R12019         05 WK-OLI-ELEMENTO-NW    OCCURS 400.
R12117            07 WK-NW-PGPF-BANK-ACC-TYP  PIC X(15).
R12117            07 WK-NW-PGPF-PAYMT-TYPE    PIC X(03).
R12117            07 WK-NW-FLAG-TIPO-POS      PIC X(01).
R12117            07 WK-NW-CAUSALE-ADD        PIC X(10).
R12117            07 WK-NW-CAUSALE-ACC        PIC X(10).
R12117            07 WK-NW-CAUSALE-COM        PIC X(10).
R12117            07 WK-NW-CAUSALE-COM-ACC    PIC X(10).
R12117            07 WK-NW-CAUSALE-ADD-P      PIC X(10).
R12117            07 WK-NW-CAUSALE-ACC-P      PIC X(10).
R12117            07 WK-NW-CAUSALE-COM-P      PIC X(10).
R12117            07 WK-NW-CAUSALE-COM-ACC-P  PIC X(10).
R12117     03 WK-FLAG-TIPO-POS              PIC X(01).
           03 WK-TS-REC.
               05  SSAA                     PIC X(4).
               05  MM                       PIC X(2).
               05  GG                       PIC X(2).
               05  HH                       PIC X(2).
               05  MI                       PIC X(2).
               05  SS                       PIC X(2).
               05  NNNNNN                   PIC X(6).
           03  WK-TS-HV.
               05  SSAA                     PIC X(4).
               05  FILLER                   PIC X VALUE '-'.
               05  MM                       PIC X(2).
               05  FILLER                   PIC X VALUE '-'.
               05  GG                       PIC X(2).
               05  FILLER                   PIC X VALUE '-'.
               05  HH                       PIC X(2).
               05  FILLER                   PIC X VALUE '.'.
               05  MI                       PIC X(2).
               05  FILLER                   PIC X VALUE '.'.
               05  SS                       PIC X(2).
               05  FILLER                   PIC X VALUE '.'.
               05  NNNNNN                   PIC X(6).
      *
           03  COMSD50-FILIALE              PIC X(05).
           03  COMSD50-CATRAPP              PIC X(04).
      *
R14316     03  WS-AREA-APPO-YPOE-DESC       PIC X(0035).
R14316*
           03  WS-AREA-APPO-YPOE.
R14316*        05  WS-AREA-APPO-YPOE-DESC   PIC X(0035).
R14316         05  WS-AREA-APPO-YPOE-D      PIC X(0035).
               05  FILLER                   PIC X(0001).
R14316*        05  WS-AREA-APPO-PAYEMT-UID  PIC X(0018).
R11422*        05  WS-AREA-APPO-YPOE-IBAN   PIC X(0031).
R11422         05  WS-AREA-APPO-YPOE-NRFA   PIC X(0012).
R11422         05  FILLER                   PIC X(0001).
R11422         05  WS-AREA-APPO-YPOE-MID    PIC X(0012).
               05  FILLER                   PIC X(0001).
R14316*        05  WS-AREA-APPO-ACC-TYPE    PIC X(0020).
R14316         05  WS-AREA-APPO-YPOE-IMPO   PIC X(0012).
      *
           03  WS-PRIM-VOLT                 PIC X(01).
               88  WS-PRIM-VOLT-SI          VALUE '1'.
               88  WS-PRIM-VOLT-NO          VALUE '0'.
      *
           03  W100-CAUSALE                 PIC X(10).
           03  W100-CODOPE                  PIC X(10).
           03  W100-BANK-ACC-TYPE           PIC X(10).
      *
           03  YP-MSGERR.
               05  YP-MSGERR-1              PIC X(80).
               05  YP-MSGERR-2              PIC X(80).
               05  YP-MSGERR-3              PIC X(80).
               05  YP-MSGERR-4              PIC X(80).
               05  YP-MSGERR-5              PIC X(80).
           03  FILLER REDEFINES YP-MSGERR.
               05  YP-MSGOCC  OCCURS 5 TIMES PIC X(80).
      *
           03  COM-ID-REC                   PIC X(07).
               88  REC-DI-TESTA             VALUE '0610697'.
               88  REC-DI-CODA              VALUE '0610695'.
               88  REC-DATI                 VALUE '0610200','0610300'.
           03  COM-CURR-DATE                PIC X(10).
           03  P303-FILE-STATUS             PIC XX.
               88  RC-NORMAL                VALUE '00' '97' '  '.
               88  RC-EOF                   VALUE '10'.
               88  RC-NOTFND                VALUE '23'.
               88  RC-DUPKEY                VALUE '22'.
           03  WS-DATA-GGMMSSAA             PIC X(10).
           03  WS-ORA-DAY-2                 PIC X(12).
           03  WS-DATA-AAAAMMGG             PIC X(08).

R15420*  ---- Area passaggio dati per pgm Z3UCUI99
R15420     COPY Z3CLUI99.
R15420
R14316*  ---- AREA UTILIZZATA PER RICHIAMARE MODULI DI CARD
R14316     COPY Z3CWDCOM  REPLACING =='Z3CWDCOM'== BY ==Z3CWDCOM==.
R14316
R14316*  ---- COPY GESTIONE IBAN
R14316     COPY Z3CLSPO2.
R14316
R14316*  ---- COPY X RECUPERO PAN II, SALDO E CAPACITA' NOMINALE CARTA
R14316     COPY Z3CLUIFA.
R14316
R14316*  ---- NOMI ROUTINE BATCH
R14316     COPY Z3CWNORB.
R14316
R14316*  ---- AREA PASSAGGIO DATI PER PGM Z3BCGE90
R14316     COPY Z3CLGE90.
R14316
R12019*  ---- TRACCIATO FLUSSO OUT BILL CCB
R12019     COPY  YPCRBILC.
R12019
TK1274*---------------------------------------------------------------
TK1274*   TABELLA GEP FPR - Param.Contabili per Coec di Pos Pagopa
TK1274*   monoente
TK1274*---------------------------------------------------------------
TK1274     COPY YPCRTFPR REPLACING 'YPCRTFPR' BY YPCRTFPR
TK1274                                 'TFPR' BY     TFPR.
      *---------------------------------------------------------------
      *   TABELLA GEP FAD CONTENENTE PER TIPOLOGIA RAPPORTO, ALCUNI
      *   DATI RELATIVI ALLA COMPILAZIONE DEL TRACCIATO D50
      *---------------------------------------------------------------
R11422     COPY YPCRTFAD REPLACING 'YPCRTFAD' BY YPCRTF03
R11422                                 'TFAD' BY     TF03.
R14316*  ---- TRACCIATO FLUSSO CONTABILE NORMALIZZATO
R14316     COPY  XYCRCONT.
R14316
R14316*  ---- TRACCIATO RECORD DCD
R14316 01 FILLER.
R14316     COPY VVG0000R.
R14316
R14316*  ---- DECODIFICA STRINGA DECIMALE IN ESADECIMALE E VICEVERSA
R14316 01 YPCCP008.
R14316    COPY YPCCP008 REPLACING ==(PREFIX)== BY ==YPCCP008==.
R14316
R05818*  ----  COPY PER IL REPERIMENTO DEI DATI ANAGRAFICI
R05818 01 AREA-ACS108A.
R05818     COPY ACS108A.
      *  ---- TRACCIATO RECORD DELLA TABELLA 'PI '
R15716*    COPY YYCRTPI.
R15716*
R15716 01  FILLER                           PIC X(08) VALUE 'YPCRTPI'.
R15716     COPY YPCRTPI
R15716          REPLACING =='YPCRTPI'== BY ==YPCRTPI==.
R15716 01  FILLER                           PIC X(08) VALUE 'WSCRTPI'.
R15716     COPY YPCRTPI
R15716          REPLACING =='YPCRTPI'== BY ==WSCRTPI==.
R12117 01  FILLER                           PIC X(08) VALUE 'YPCWFAMI'.
R12117    COPY YPCWFAMI.

      *  ---- COPY PER ROUTINE XYRC0005 GESTIONE ABEND.
           COPY  XYCW0005.
      *  ---- AREA PASSAGGIO DATI PER PGM YPBCREQD
R11422     COPY YPCRREQX.


      *  ---- TRACCIATO FILE PGPF INPUT
           COPY  YPCRPGPF.

      *  ---- TRACCIATO FILE OPECONT
           COPY  CRVFSD50.

      * ----- COMMAREA MODULO DI INCCIRY CONTO CORRENTE
           COPY PPTCINCC.
      *
R05316* ----- TRACCIATO RECORD DELLA TABELLA 'PGP'
R05316     COPY YPCRTPGP.
R05818* ----- TRACCIATO RECORD DELLA TABELLA 'CCB'
R05818     COPY YPCRTCCB.
      *
      *--* Tracciati per record di file output
      *
       COPY  YYCWUTDA.
R04417*
R04417* ----- TRACCIATO PER LA CALL ALLA ROUTINE DELLA DATA
R04417     COPY XSADAT.
R04417 01 WS-XSCDAT                        PIC X(08) VALUE 'XSCDAT'.
      *
R12117 01 YPRCFAMI                         PIC X(08) VALUE 'YPRCFAMI'.
      *-------VARIABILI DI COMODO PER CALCOLO DATA
R11422 01  WS-DATA                              PIC 9(6).
R11422 01  W100-DATA-CHIUSURA-MM-1           PIC X(08) VALUE SPACES.
R11422 01  W100-DATA-CHIUSURA-MM-2           PIC X(08) VALUE SPACES.
R11422 01  WS-DATA-8                            PIC X(8).
      *
       01  P303-MSG-ERRORE.
           02  P303-MSG-ERRORE-1.
           03  FILLER                      PIC X(11)  VALUE
               'PROGRAMMA: '.
           03  P303-MSGER-PGM              PIC X(08)  VALUE SPACES.
           03  FILLER                      PIC X(12)  VALUE
               '  . RIFER. ('.
           03  P303-MSGER-RIF              PIC X(02)  VALUE ZERO.
           03  FILLER                      PIC X(04)  VALUE
               ') - '.
           03  P303-MSGER-TIPO             PIC X(08)  VALUE SPACES.
           03  FILLER                      PIC X(03)  VALUE ' - '.
           03  P303-MSGER-FILE             PIC X(08)  VALUE SPACES.
           03  FILLER                      PIC X(11)  VALUE
               ' - STATUS: '.
           03  P303-MSGER-STATUS           PIC X(06)  VALUE SPACES.
           03  FILLER                      PIC X(07)  VALUE ' DATO: '.
           03  P303-MSGER-DATO             PIC X(20)  VALUE SPACES.
           03  FILLER                      PIC X(04)  VALUE ' -- '.
           03  P303-MSGER-DESCR            PIC X(30)  VALUE SPACES.
           02  P303-MSG-ERRORE-2           PIC X(50)  VALUE SPACES.
      *----------------------------------------------------------------
******* ATTENZIONE   DESCRIZIONE CONTI CORRENTI
*******              CAMPO CRVSD50-DESCMOV MAX 120 BYTE
       01  W100-DESCMOV.
R05818     03  FILLER                 PIC X(30) VALUE
FF0518*    'Full Acquiring Poste Italiane '.
FF0518     'Full Acquiring PostePay S.p.A.'.
           03  FILLER                 PIC X(05) VALUE
           ' del '.
           03  W100-DATA-OPERAZ       PIC X(10) VALUE SPACES.
160408     03  FILLER                 PIC X(05) VALUE ' Mid '.
160408     03  W100-MERCHANT          PIC X(15).
170511     03  FILLER                 PIC X(01) VALUE SPACE.
170511     03  W100-NICKNAME          PIC X(40).
           03  W100-DESC-BRAND        PIC X(07) VALUE ' Brand '.
           03  W100-BRAND-CODE        PIC X(04).
           03  FILLER                 PIC X(03) VALUE SPACE.
      *
       01  W100-DESCMOV-200.
           03  FILLER                 PIC X(01) VALUE SPACES.
           03  FILLER                 PIC X(30) VALUE
FF0518*    'Full Acquiring Poste Italiane'.
FF0518     'Full Acquiring PostePay S.p.A.'.
160408     03  FILLER                 PIC X(12) VALUE ' - Merchant '.
160408     03  W100-MERCHANT-CH       PIC X(15).
170511     03  FILLER                 PIC X(01) VALUE SPACE.
170511     03  W100-NICKNAME-CH       PIC X(40).
           03  W100-DESC-BRAND-CH     PIC X(07) VALUE ' Brand '.
           03  W100-BRAND-CODE-CH     PIC X(04).
      *
R05818 01  W100-DESCMOV-CC-VTPIE.
R05818     03  FILLER                 PIC X(01) VALUE SPACES.
R05818     03  FILLER                 PIC X(30) VALUE
FF0518*    'Full Acquiring Poste Italiane'.
FF0518     'Full Acquiring PostePay S.p.A.'.
R05818     03  FILLER                 PIC X(16) VALUE
R05818     ' Operazioni del '.
R05818     03  W100-DATA-OPER         PIC X(10) VALUE SPACES.
R05818     03  FILLER                 PIC X(25) VALUE
R05818     ' Incassi su IBAN Bancario'.
      *
R14316 01  W100-DESPPAY.
R14316     03  FILLER                 PIC X(15) VALUE
R14316     'Operazioni del '.
R14316     03  W100-DESPPAY-DT        PIC X(10) VALUE SPACES.
R14316     03  FILLER                 PIC X(10) VALUE ' Merchant '.
R14316     03  W100-DESPPAY-ME        PIC X(15).
R14316     03  W100-DESPPAY-DESC-BR   PIC X(07) VALUE ' Brand '.
R14316     03  W100-DESPPAY-BR        PIC X(04).
R14316*
      *
R14316 01  W100-DESPPAY-200.
R14316     03  FILLER                 PIC X(09) VALUE 'Merchant '.
R14316     03  W100-MERCHANT-PPAY     PIC X(15).
R14316     03  W100-DESC-BRAND-PPAY   PIC X(07) VALUE ' Brand '.
R14316     03  W100-BRAND-CODE-PPAY   PIC X(04).
      *
       01  HV-IMPORTO                    PIC S9(12)V COMP-3 VALUE ZERO.
       01  WK-NULL-ELAB                  PIC S9(04)  COMP.
       01  WK-NULL-VAR                   PIC S9(04)  COMP.
       01  HV-FUNCTION-CODE              PIC  X(03).
       01  HV-SUMMARY-UID                PIC  X(18).
       01  HV-PAYMENT-UID                PIC  X(18).
      *--------------------------------------------------------
       01 HV-TABE.
          05 HV-TABE-KNAMTB1                PIC X(0003).
          05 HV-TABE-KVARTB1                PIC X(0027).
          05 HV-TABE-DATI.
             49 HV-TABE-DATI-L              PIC  S9(4) COMP.
             49 HV-TABE-DATI-A              PIC X(2000).
R14217*01 HV-KEY-RANDOM           PIC X(08).
R10418 01 HV-KEY-RANDOM           PIC S9(15) COMP-3.
      *================================================================*
      *    Area DB2                                                    *
      *================================================================*
      *
           EXEC SQL INCLUDE SQLCA END-EXEC.
      *
           EXEC SQL INCLUDE YPDCPGPF  END-EXEC.
R11422     EXEC SQL INCLUDE YPDCFAS2  END-EXEC.
      *
       01  W100-APPO-SQLCODE                 PIC -----9.
      *
      *----------------------------------------------------------------
      * COPY PER SEGNALAZIONI SU SYSOUT
      *----------------------------------------------------------------
      *01  FILLER                           PIC X(08) VALUE 'XYCP303'.
      *    COPY XYCWSGNL.
      *
      *----------------------------------------------------------------
      * COPY INSTALLAZIONE
      *----------------------------------------------------------------
       01  FILLER                           PIC X(08) VALUE 'XYCW0005'.
           COPY XYCW0005.
      *
      * ----- UTILITY PER STAMPA ELAB/CALCOLO DATE/ERRORI VSAM
       COPY YPCWS001.
      *
      *================================================================*
       PROCEDURE DIVISION.
      *
      *--* Label iniziali
           PERFORM INIZ-WORK              THRU F-INIZ-WORK
           PERFORM STAM-PRIM-RIGH         THRU F-STAM-PRIM-RIGH
           PERFORM OPEN-FILE              THRU F-OPEN-FILE
           PERFORM LEGG-FILE              THRU F-LEGG-FILE
      *
      *--* Controlli iniziali
           PERFORM CTRL-INIZ              THRU F-CTRL-INIZ
      *
R11422*--* Gestione data da routine
R11422     PERFORM ROUT-DATE            THRU F-ROUT-DATE
R14316*
R14316*--* Legge la tabella YYDTABE ("PI/OLI/AESMACPREP")
R12117*    PERFORM LEGG-TABE-PREP         THRU F-LEGG-TABE-PREP
      *
      *--* Legge la tabella YYDTABE ("PI/OLI/AESMAC")
R12117*    PERFORM LEGG-TABE              THRU F-LEGG-TABE
R12117*--* Legge la tabella YYDTABE ("PI/OLI/SMAC")
R12117     PERFORM START-TABE-SMAC        THRU F-START-TABE-SMAC
R12117     PERFORM ELAB-TABE-SMAC         THRU F-ELAB-TABE-SMAC
R12117          UNTIL FINE-TABE-SMAC
R12019*         OR IND1 = 200.
R12019          OR IND1 = 400.
      *
      *--* Elaborazione principale
           PERFORM ELAB                   THRU F-ELAB
                                          UNTIL WS-EOF-IPAYMENT = 1
R05818     IF WS-NUM-OPER-CC  > ZERO
R05818         PERFORM IMPOSTA-DATI-OUT-CC-VTPIE
R05818            THRU EX-IMPOSTA-DATI-OUT-CC-VTPIE
R05818         PERFORM SCRIVI-REC-OUT
R05818            THRU EX-SCRIVI-REC-OUT
R05818         PERFORM IMPOSTA-DATI-T-FXML
R05818             THRU EX-IMPOSTA-DATI-T-FXML
R05818         PERFORM SCRIVI-FXML
R05818             THRU EX-SCRIVI-FXML
R05818     END-IF

R11422     IF WS-NUM-OPER-CC-B  > ZERO
R11422         PERFORM IMPOSTA-DATI-T-FXM2
R11422             THRU EX-IMPOSTA-DATI-T-FXM2
R11422         PERFORM SCRIVI-FXM2
R11422             THRU EX-SCRIVI-FXM2
R11422     END-IF
      *
      *--* Controlli finali
           PERFORM CTRL-FINA              THRU F-CTRL-FINA
      *
           PERFORM STAM-RIGH-TOTA         THRU F-STAM-RIGH-TOTA
           PERFORM STAM-RIGH-FINA         THRU F-STAM-RIGH-FINA
           .
       FINE-JOB.
           PERFORM CLOSE-FILE             THRU F-CLOSE-FILE
           GOBACK
           .
      *================================================================*
       INIZ-WORK.
      *
           MOVE ZERO                        TO WS-ZERO
      *
           MOVE SPACE                       TO WS-SPACE
      *
      *--* Setta parametri
           SET WS-PRIM-VOLT-SI              TO TRUE
R05818     SET WS-LETTA-GEP-CCB-NO          TO TRUE
      *
      *--* Acquisizione data
           PERFORM IMPO-DATA              THRU EX-IMPO-DATA
           .
       F-INIZ-WORK.
           EXIT.
      *==============================================================*
       IMPO-DATA.
      *
           ACCEPT WS-DATA-AAMMGG          FROM DATE
      *--* Formato SSAAMMGG
           MOVE '20'                        TO WS-DATA-AAAAMMGG(1:2)
010715     MOVE WS-DATA-AAMMGG              TO WS-DATA-AAAAMMGG(3:6)
      *
      *--* Formato GGMMSSAA
           STRING WS-DATA-AAMMGG(5:2) '-'
                  WS-DATA-AAMMGG(3:2) '-'
                  '20' WS-DATA-AAMMGG(1:2)
           DELIMITED BY SIZE              INTO WS-DATA-GGMMSSAA
           END-STRING
           ACCEPT WS-ORA                  FROM TIME
           .
       EX-IMPO-DATA.
           EXIT.
      *================================================================*
      *    Routine di segnalazione inizio programma                    *
      *================================================================*
       STAM-PRIM-RIGH.
      *
           COPY  YPCPS001.
      *
           MOVE  SPACES                     TO YPCWS001-TEST-1
      *
           STRING '* FULL ACQUIRING * '
                  'CARICA LA TAB.YPTBPGPF E CREA FLUSSO CON D50 X POZZO'
                  '- PROGRAMMA: YPBCEPGP -'
           DELIMITED BY SIZE              INTO YPCWS001-TEST-1
           END-STRING
      *
           MOVE SPACES                      TO YPCWS001-RIGA
           PERFORM SCRIVI-ST              THRU EX-SCRIVI-ST
           MOVE SPACES                      TO YPCWS001-RIGA
           .
       F-STAM-PRIM-RIGH.
           EXIT.
      *================================================================*
      *    Apertura archivi                                            *
      *================================================================*
       OPEN-FILE.
      *
           OPEN INPUT  IPAYMENT.
           IF NOT IPAY-NORMAL
              MOVE ST-IPAYMENT           TO P303-FILE-STATUS
              MOVE '01'                  TO P303-MSGER-RIF
              MOVE 'OPENINP'             TO P303-MSGER-TIPO
              MOVE 'IPAYMENT'            TO P303-MSGER-FILE
              PERFORM ERRORE-P303        THRU EX-ERRORE-P303
           END-IF.

           OPEN INPUT  YYDTABE
           IF NOT YYDT-NORMAL
              MOVE ST-YYDTABE            TO P303-FILE-STATUS
              MOVE '02'                  TO P303-MSGER-RIF
              MOVE 'OPENINP'             TO P303-MSGER-TIPO
              MOVE 'YYDTABE '            TO P303-MSGER-FILE
              PERFORM ERRORE-P303       THRU EX-ERRORE-P303
           END-IF.

           OPEN OUTPUT OPECONT
           IF NOT OPEC-NORMAL
              MOVE ST-OPECONT            TO P303-FILE-STATUS
              MOVE '03'                  TO P303-MSGER-RIF
              MOVE 'OPENOUT'             TO P303-MSGER-TIPO
              MOVE 'OPECONT'             TO P303-MSGER-FILE
              PERFORM ERRORE-P303       THRU EX-ERRORE-P303
           END-IF.
FIANNH
FIANNH     OPEN OUTPUT OPECONTB
FIANNH     IF NOT OPECB-NORMAL
FIANNH        MOVE ST-OPECONTB           TO P303-FILE-STATUS
FIANNH        MOVE '03'                  TO P303-MSGER-RIF
FIANNH        MOVE 'OPENOUT'             TO P303-MSGER-TIPO
FIANNH        MOVE 'OPECONTB'             TO P303-MSGER-FILE
FIANNH        PERFORM ERRORE-P303       THRU EX-ERRORE-P303
FIANNH     END-IF.
FIANNH
R14316
R14316     OPEN OUTPUT XYDCONT
R14316     IF NOT OPECB-NORMAL
R14316        MOVE ST-XYDCONT            TO P303-FILE-STATUS
R14316        MOVE '16'                  TO P303-MSGER-RIF
R14316        MOVE 'OPENOUT'             TO P303-MSGER-TIPO
R14316        MOVE 'XYDCONT'              TO P303-MSGER-FILE
R14316        PERFORM ERRORE-P303       THRU EX-ERRORE-P303
R14316     END-IF.
R14316
R14316     OPEN OUTPUT OUTDCD
R14316     IF NOT OPECB-NORMAL
R14316        MOVE ST-OUTDCD             TO P303-FILE-STATUS
R14316        MOVE '18'                  TO P303-MSGER-RIF
R14316        MOVE 'OPENOUT'             TO P303-MSGER-TIPO
R14316        MOVE 'OUTDCD'               TO P303-MSGER-FILE
R14316        PERFORM ERRORE-P303       THRU EX-ERRORE-P303
R14316     END-IF.

           OPEN OUTPUT OUSCARTI.
           IF NOT OUSC-NORMAL
              MOVE ST-OUSCARTI           TO P303-FILE-STATUS
              MOVE '04'                  TO P303-MSGER-RIF
              MOVE 'OPENOUT'             TO P303-MSGER-TIPO
              MOVE 'OUSCARTI'            TO P303-MSGER-FILE
              PERFORM ERRORE-P303        THRU EX-ERRORE-P303
           END-IF.

R11817     OPEN OUTPUT OUSCART2.
R11817     IF NOT OUS2-NORMAL
R11817        MOVE ST-OUSCART2           TO P303-FILE-STATUS
R11817        MOVE '04'                  TO P303-MSGER-RIF
R11817        MOVE 'OPENOUT'             TO P303-MSGER-TIPO
R11817        MOVE 'OUSCART2'            TO P303-MSGER-FILE
R11817        PERFORM ERRORE-P303        THRU EX-ERRORE-P303
R11817     END-IF.

R05818     OPEN OUTPUT OUTFXML.
R05818     IF NOT OXML-NORMAL
R05818        MOVE ST-OUTFXML            TO P303-FILE-STATUS
R05818        MOVE '05'                  TO P303-MSGER-RIF
R05818        MOVE 'OPENOUT'             TO P303-MSGER-TIPO
R05818        MOVE 'OUTFXML'             TO P303-MSGER-FILE
R05818        PERFORM ERRORE-P303        THRU EX-ERRORE-P303
R05818     END-IF.

R11422     OPEN OUTPUT OUTFXM2.
R11422     IF NOT OXM2-NORMAL
R11422        MOVE ST-OUTFXM2            TO P303-FILE-STATUS
R11422        MOVE '10'                  TO P303-MSGER-RIF
R11422        MOVE 'OPENOUT'             TO P303-MSGER-TIPO
R11422        MOVE 'OUTFXM2'             TO P303-MSGER-FILE
R11422        PERFORM ERRORE-P303        THRU EX-ERRORE-P303
R11422     END-IF.

R12019     OPEN OUTPUT BILLCCB.
R12019     IF NOT BILC-NORMAL
R12019        MOVE ST-BILLCCB            TO P303-FILE-STATUS
R12019        MOVE '13'                  TO P303-MSGER-RIF
R12019        MOVE 'OPENOUT'             TO P303-MSGER-TIPO
R12019        MOVE 'BILLCCB'             TO P303-MSGER-FILE
R12019        PERFORM ERRORE-P303        THRU EX-ERRORE-P303
R12019     END-IF.

           OPEN OUTPUT YPOERRO.
           IF NOT OERR-NORMAL
R03817*       MOVE ST-OUSCARTI           TO P303-FILE-STATUS
R03817        MOVE ST-YPOERRO            TO P303-FILE-STATUS
              MOVE '14'                  TO P303-MSGER-RIF
              MOVE 'OPENOUT'             TO P303-MSGER-TIPO
              MOVE 'YPOERRO'             TO P303-MSGER-FILE
              PERFORM ERRORE-P303        THRU EX-ERRORE-P303
           END-IF.
           .
       F-OPEN-FILE.
           EXIT.
      *================================================================*
      *    Lettura record di ACQUINP                                   *
      *================================================================*
       LEGG-FILE.
      *
           READ IPAYMENT                  INTO YPCRPGPF-WORK
      *
           IF IPAY-EOF
              MOVE 1 TO WS-EOF-IPAYMENT
              GO TO F-LEGG-FILE
           END-IF
      *
           IF NOT IPAY-NORMAL
              MOVE ST-IPAYMENT              TO P303-FILE-STATUS
              MOVE '10'                     TO P303-MSGER-RIF
              MOVE 'READ'                   TO P303-MSGER-TIPO
              MOVE 'IPAYMENT'               TO P303-MSGER-FILE
              MOVE SPACES                   TO P303-MSGER-DESCR
              PERFORM ERRORE-P303         THRU EX-ERRORE-P303
           END-IF
      *
           MOVE PGPFH-FUNCT-CODE            TO WS-TIPO-FUN-CODE
           ADD  1                           TO CTR-CONT-LETTI-TOT
           .
       F-LEGG-FILE.
           EXIT.
      *================================================================*
       CTRL-INIZ.
      *
      *--* Controlla se il file � vuoto
           IF WS-EOF-IPAYMENT = 1
              PERFORM   IMPO-ERRO-FILE-VUOTO
                 THRU F-IMPO-ERRO-FILE-VUOTO
           END-IF
      *
      *--* Controlla che il primo record sia un record di testa
           IF  PGPFH-MSG-TYPE-ID = '0610'
           AND PGPFH-FUNCT-CODE  = '697'
              CONTINUE
           ELSE
              PERFORM   IMPO-ERRO-PRIM-RECO
                 THRU F-IMPO-ERRO-PRIM-RECO
           END-IF
           .
       F-CTRL-INIZ.
           EXIT.
      *================================================================*
R14316 LEGG-TABE-PREP.
      *
           INITIALIZE                          WSCRTPI
      *
           MOVE 'PI '                       TO WSCRTPI-COD
           MOVE 'OLI'                       TO WSCRTPI-COD-VAR
           MOVE 'AESMACPREP'                TO WSCRTPI-KEY-TAB
      *
           MOVE WSCRTPI-KEY                 TO YYDTABE-KEY

           READ YYDTABE                   INTO WSCRTPI
      *
           IF NOT YYDT-NORMAL
              MOVE ST-YYDTABE               TO P303-FILE-STATUS
              MOVE 'READTAB'                TO P303-MSGER-TIPO
              MOVE 'YYDTABE'                TO P303-MSGER-FILE
              MOVE YYDTABE-KEY              TO P303-MSGER-DATO
              IF YYDT-NOTFND
                 MOVE '11'                  TO P303-MSGER-RIF
                 MOVE 'ELEM."PREP" NON TROVATO'
                                            TO P303-MSGER-DESCR
               ELSE
                 MOVE '12'                  TO P303-MSGER-RIF
                 MOVE 'ERRORE IN LETTURA ELE."PREP"'
                                            TO P303-MSGER-DESCR
              END-IF
              MOVE 'JOB INTERROTTO'         TO P303-MSG-ERRORE-2
              PERFORM ERRORE-P303         THRU EX-ERRORE-P303
           END-IF
           .
       F-LEGG-TABE-PREP.
           EXIT.
      *================================================================*
R12117 START-TABE-SMAC.
      *
           MOVE SPACE   TO WK-TABE-SMAC
           INITIALIZE                          YPCRTPI
      *
           MOVE 'PI '                       TO YPCRTPI-COD
           MOVE 'OLI'                       TO YPCRTPI-COD-VAR
           MOVE 'SMAC'                      TO YPCRTPI-KEY-TAB
      *
           MOVE YPCRTPI-KEY                 TO YYDTABE-KEY
           START YYDTABE
                 KEY GREATER THAN YYDTABE-KEY
      *
           IF NOT YYDT-NORMAL
              MOVE ST-YYDTABE               TO P303-FILE-STATUS
              MOVE 'START  '                TO P303-MSGER-TIPO
              MOVE 'YYDTABE'                TO P303-MSGER-FILE
              MOVE YYDTABE-KEY              TO P303-MSGER-DATO
              IF YYDT-NOTFND
                 MOVE '06'                   TO P303-MSGER-RIF
                 MOVE 'ELEM."SMAC" NON TROVATO'
                                             TO P303-MSGER-DESCR
               ELSE
                 MOVE '07'                   TO P303-MSGER-RIF
                 MOVE 'ERRORE IN LETTURA ELE."SMAC"'
                                             TO P303-MSGER-DESCR
              END-IF
              MOVE 'JOB INTERROTTO'         TO P303-MSG-ERRORE-2
              PERFORM ERRORE-P303         THRU EX-ERRORE-P303
           END-IF
           .
R12117 F-START-TABE-SMAC.
           EXIT.
      *================================================================*
R12117 ELAB-TABE-SMAC.
      *
           PERFORM LEGG-TABE-SMAC
              THRU F-LEGG-TABE-SMAC
           IF NOT FINE-TABE-SMAC
             ADD 1 TO IND1
             PERFORM CARI-TABE-SMAC
                THRU F-CARI-TABE-SMAC
           END-IF
           .
R12117 F-ELAB-TABE-SMAC.
           EXIT.
      *================================================================*
R12117 LEGG-TABE-SMAC.
      *
           MOVE YPCRTPI-KEY                 TO YYDTABE-KEY

           READ  YYDTABE NEXT.
           MOVE YYDTABE-REC-FD              TO YPCRTPI
      *
           IF NOT YYDT-NORMAL
              MOVE ST-YYDTABE               TO P303-FILE-STATUS
              MOVE 'READTAB'                TO P303-MSGER-TIPO
              MOVE 'YYDTABE'                TO P303-MSGER-FILE
              MOVE YYDTABE-KEY              TO P303-MSGER-DATO
              IF YYDT-NOTFND
                 SET  FINE-TABE-SMAC  TO TRUE
               ELSE
                 MOVE '07'                   TO P303-MSGER-RIF
                 MOVE 'ERRORE IN LETTURA ELE."SMAC"'
                                             TO P303-MSGER-DESCR
              END-IF
              MOVE 'JOB INTERROTTO'         TO P303-MSG-ERRORE-2
              PERFORM ERRORE-P303         THRU EX-ERRORE-P303
           END-IF
           IF YPCRTPI-KEY-SMAC NOT = 'SMAC'
               SET  FINE-TABE-SMAC  TO TRUE
           END-IF
           .
R12117 F-LEGG-TABE-SMAC.
           EXIT.
      *================================================================*
R12117 CARI-TABE-SMAC.
      *
           MOVE YPCRTPI-OLI-NW-CAU-ADD
                                     TO WK-NW-CAUSALE-ADD(IND1)
           MOVE YPCRTPI-OLI-NW-CAU-ACC
                                     TO WK-NW-CAUSALE-ACC(IND1)
           MOVE YPCRTPI-OLI-NW-CAU-COM
                                     TO WK-NW-CAUSALE-COM(IND1)
           MOVE YPCRTPI-OLI-NW-CAU-COM-ACC
                                   TO WK-NW-CAUSALE-COM-ACC(IND1)
           MOVE YPCRTPI-OLI-NW-CAU-ADD-P
                                     TO WK-NW-CAUSALE-ADD-P(IND1)
           MOVE YPCRTPI-OLI-NW-CAU-ACC-P
                                     TO WK-NW-CAUSALE-ACC-P(IND1)
           MOVE YPCRTPI-OLI-NW-CAU-COM-P
                                     TO WK-NW-CAUSALE-COM-P(IND1)
           MOVE YPCRTPI-OLI-NW-CAU-COM-ACC-P
                                   TO WK-NW-CAUSALE-COM-ACC-P(IND1)
           MOVE YPCRTPI-PGPF-BANK-ACC-TYP
                                 TO WK-NW-PGPF-BANK-ACC-TYP(IND1)
           MOVE YPCRTPI-PGPF-PAYMT-TYPE
                                 TO WK-NW-PGPF-PAYMT-TYPE(IND1)
           MOVE YPCRTPI-FLAG-TIPO-POS
                                 TO WK-NW-FLAG-TIPO-POS(IND1)
           .
R12117 F-CARI-TABE-SMAC.
           EXIT.
      *================================================================*
       LEGG-TABE.
      *
           INITIALIZE                          YPCRTPI
      *
           MOVE 'PI '                       TO YPCRTPI-COD
           MOVE 'OLI'                       TO YPCRTPI-COD-VAR
           MOVE 'AESMAC'                    TO YPCRTPI-KEY-TAB
      *
           MOVE YPCRTPI-KEY                 TO YYDTABE-KEY

           READ YYDTABE                   INTO YPCRTPI
      *
           IF NOT YYDT-NORMAL
              MOVE ST-YYDTABE               TO P303-FILE-STATUS
              MOVE 'READTAB'                TO P303-MSGER-TIPO
              MOVE 'YYDTABE'                TO P303-MSGER-FILE
              MOVE YYDTABE-KEY              TO P303-MSGER-DATO
              IF YYDT-NOTFND
                 MOVE '06'                   TO P303-MSGER-RIF
                 MOVE 'ELEMENTO NON TROVATO' TO P303-MSGER-DESCR
               ELSE
                 MOVE '07'                   TO P303-MSGER-RIF
                 MOVE 'ERRORE IN LETTURA   ' TO P303-MSGER-DESCR
              END-IF
              MOVE 'JOB INTERROTTO'         TO P303-MSG-ERRORE-2
              PERFORM ERRORE-P303         THRU EX-ERRORE-P303
           END-IF
           .
       F-LEGG-TABE.
           EXIT.
      *================================================================
       ELAB.
DBG==>*    DISPLAY '--ELABORA===================================='
      *-
           SET WS-SALT-ELAB-NO              TO TRUE
R05316     SET WS-SALT-CONT-NO              TO TRUE
R00317     SET WS-TIPO-BILL-NO              TO TRUE
           SET WS-SCRI-SCAR-NO              TO TRUE
R07420     SET WS-CONTAB-NO                 TO TRUE
R03817     SET WS-IMPO-ZERO-NO              TO TRUE
R11422     SET WS-ELAB-CC-BANCARI-NO        TO TRUE
      *-
      *--* Esegue controlli sul record letto
           PERFORM CTRL-MSG
           THRU  F-CTRL-MSG
DBG==>*    DISPLAY 'WS-SALT-ELAB     -1-  ('WS-SALT-ELAB')'
           IF WS-SALT-ELAB-NO
      *-
      *--* Acquisce dati
              PERFORM ACQU-DATI
              THRU  F-ACQU-DATI
DBG==>*    DISPLAY 'DOPO ACQU-DATI '
DBG==>*    DISPLAY 'WS-SALT-ELAB     -2-  ('WS-SALT-ELAB')'
              IF WS-SALT-ELAB-NO
      *-
      *--* Inserisce dati sulla tabella PGPF
TEST             PERFORM INSE-TABE-PGPF
TEST             THRU  F-INSE-TABE-PGPF
                 IF WS-SALT-ELAB-NO
R05316          AND WS-SALT-CONT-NO
R08421          AND PGPF-FUNCT-CODE NOT = '301'
      *-
      *--* Imposta dati per scrivere file di output
R14316*--* Se siamo nel caso di Postepay Evolution Business
R14316*--* imposta i campi della copy XYCRCONT altrimenti
R14316*--* negli altri casi continua ad impostare e scrivere un D50
DBG==>*      display 'WS-PAN-PPAY-EVOL-BUSI : ' WS-PAN-PPAY-EVOL-BUSI
R14316              IF WS-PAN-PPAY-EVOL-BUSI-SI
R14316                 PERFORM IMPOSTA-DATI-OUT-CONT
R14316                 THRU EX-IMPOSTA-DATI-OUT-CONT
R14316                 PERFORM SCRIVI-REC-OUT-CONT
R14316                 THRU EX-SCRIVI-REC-OUT-CONT
R14316              ELSE
DBG==>*              DISPLAY 'WS-TIPO-CC          ('WS-TIPO-CC ')'
DBG==>*              DISPLAY 'PGPF-DB-CR-FLAG     ('PGPF-DB-CR-FLAG ')'
DBG==>*              DISPLAY 'PGPF-BANK-ACC-TYP   ('PGPF-BANK-ACC-TYP')'
DBG==>*              DISPLAY 'PGPF-PAYEMT-UID     ('PGPF-PAYEMT-UID ')'
DBG==>*              DISPLAY 'YPCWFAMI-O-CONT-BILLI-CCB ('
DBG==>*                       YPCWFAMI-O-CONT-BILLI-CCB ')'
R05818               IF WS-CC-POSTE
                       PERFORM IMPOSTA-DATI-OUT
                       THRU EX-IMPOSTA-DATI-OUT
FIANNH                 EVALUATE  YPCRTPGP-DATI-TIPO-ELAB
FIANNH                 WHEN 'IBAN'
                          PERFORM SCRIVI-REC-OUT
                          THRU EX-SCRIVI-REC-OUT
FIANNH                 WHEN 'IBAB'
FIANNH                    PERFORM SCRIVI-REC-OUTB
FIANNH                    THRU EX-SCRIVI-REC-OUTB
FIANNH                 END-EVALUATE
R05818               END-IF
R05818               IF WS-CC-BANCARIO
R05818                 IF  PGPF-DB-CR-FLAG = 'D'
R12019*????                EVALUATE PGPF-PAYMT-TYPE
DBG==>*    DISPLAY 'PGPF-BANK-ACC-TYP         ('PGPF-BANK-ACC-TYP')'
DBG==>*    DISPLAY 'YPCWFAMI-O-CONT-BILLI-CCB ('
DBG==>*             YPCWFAMI-O-CONT-BILLI-CCB ')'
R12019                     EVALUATE PGPF-BANK-ACC-TYP
R12019                       WHEN 'DSC'
R12019                          IF YPCWFAMI-O-CONT-BILLI-CCB = SPACES
R12019                          OR YPCWFAMI-O-CONT-BILLI-CCB = LOW-VALUE
R12019                             SET WS-SCRI-SCAR-SI TO TRUE
R12019                          ELSE
TK1274                           PERFORM CONTROLLA-GEP-FPR
TK1274                             THRU EX-CONTROLLA-GEP-FPR
TK1274                           IF TROVATO-SU-FPR-SI
TK1274                             PERFORM IMPOSTA-BILLCCB-2
TK1274                             THRU EX-IMPOSTA-BILLCCB-2
TK1274                           ELSE
R12019                             PERFORM IMPOSTA-BILLCCB
R12019                             THRU EX-IMPOSTA-BILLCCB
R12019                             PERFORM SCRIVI-BILLCCB
R12019                             THRU EX-SCRIVI-BILLCCB
TK1274                           END-IF
R12019                          END-IF
R12019                       WHEN OTHER
ZANCHI* Se si tratta di un canone non fatturato mando a billing
ZANCHI                          IF (YPCRTPGP-DATI-TIPO-ELAB = 'IBAB'
ZANCHI*                         OR                            'IBAN'
ZANCHI                             )
ZANCHI                          AND YPCWFAMI-O-CONT-BILLI-CCB > SPACES
R12019                              PERFORM IMPOSTA-BILLCCB
R12019                                 THRU EX-IMPOSTA-BILLCCB
R12019                              PERFORM SCRIVI-BILLCCB
R12019                                 THRU EX-SCRIVI-BILLCCB
ZANCHI                          ELSE
R05818                             SET WS-SCRI-SCAR-SI TO TRUE
ZANCHI                          END-IF
R12019                     END-EVALUATE
R05818                 END-IF
R05818                 IF  PGPF-DB-CR-FLAG = 'C'
R05818                    PERFORM ADD-TOTALI
R05818                    THRU EX-ADD-TOTALI
R05818                    PERFORM IMPOSTA-DATI-FXML
R05818                    THRU EX-IMPOSTA-DATI-FXML
R05818                    PERFORM SCRIVI-FXML
R05818                    THRU EX-SCRIVI-FXML
R05818                 END-IF
R05818               END-IF
R14316              END-IF
                 END-IF
              END-IF
           END-IF
PSPSPS     .
DBG==>*    DISPLAY 'WS-SCRI-SCAR              ('WS-SCRI-SCAR     ')'
DBG==>*    DISPLAY 'PGPF-FUNCT-CODE           ('PGPF-FUNCT-CODE  ')'
      *-
      *--* In caso debba segnalare l'errore scrive lo scarto
           IF WS-SCRI-SCAR-SI
R08421     AND PGPF-FUNCT-CODE NOT = '301'
R11817*       PERFORM SCRIVI-SCARTI
R11817*       THRU EX-SCRIVI-SCARTI
R14316*--* Se siamo nel caso di Postepay Evolution Business
R11817*--* e la carta non � attiva o l'IBAN non � attivo
R14316*--* scrive un DCD in apertura
R11817*--* e il flusso scart2
R14316        IF WS-PAN-PPAY-EVOL-BUSI-SI
R03817        AND WS-IMPO-ZERO-NO
R11817        AND (   (WS-IBAN-ATTI-SI AND Z3CLUIFA-OU-RET-CODE = '002')
R11817             OR (WS-IBAN-ATTI-NO))
R14316          PERFORM SCRIVI-DCD-APER
R14316          THRU EX-SCRIVI-DCD-APER
R11817          PERFORM SCRIVI-SCART2
R11817          THRU EX-SCRIVI-SCART2
R11817        ELSE
R07420        IF  WS-CONTAB-SI
R07420          PERFORM SCRIVI-DCD-APER
R07420          THRU EX-SCRIVI-DCD-APER
R07420          PERFORM SCRIVI-SCART2
R07420          THRU EX-SCRIVI-SCART2
R07420        ELSE
R12019*--* Per volere di Zanchi � stata ripristinata la scrittura del
R12019*--* DCD ed sostituita la PERFORM SCRIVI-SCART2 con SCRIVI-SCARTI
R12019*--* solo se importo diverso da zero
R05818         IF WS-CC-BANCARIO
R12019         AND WS-IMPO-ZERO-NO
R11422          IF  (YPCWFAMI-O-ID-MANDATO   NOT = SPACES AND LOW-VALUE)
R11422          AND (PGPF-PAYEMT-UID         NOT = SPACES AND LOW-VALUE)
R11422          AND (YPCWFAMI-O-DATA-ATTIV-MANDATO
R11422                                       NOT = SPACES AND LOW-VALUE)
R11422          AND (YPCWFAMI-O-IBAN-MANDATO NOT = SPACES AND LOW-VALUE)
R11422            PERFORM TRATTA-CC-BANCARIO
R11422            THRU EX-TRATTA-CC-BANCARIO
R11422          ELSE
R13519            PERFORM SCRIVI-DCD-APER-CCB
R13519            THRU EX-SCRIVI-DCD-APER-CCB
R13519*           PERFORM SCRIVI-SCART2
R13519*           THRU EX-SCRIVI-SCART2
R13519            PERFORM SCRIVI-SCARTI
R13519            THRU EX-SCRIVI-SCARTI
DBG==>*    DISPLAY 'Imposta errore --Addebiti su IBAN Bancario--   '
R05818            MOVE 'Addebiti su IBAN Bancario           '
R05818                                      TO WS-AREA-APPO-YPOE-DESC
R11422          END-IF
R05818         ELSE
R11817          PERFORM SCRIVI-SCARTI
R11817          THRU EX-SCRIVI-SCARTI
R05818         END-IF
R07420        END-IF
R14316        END-IF
      *
      *--* Scrive file errori x mail
              PERFORM SCRIVI-ERRO-MAIL
              THRU EX-SCRIVI-ERRO-MAIL
           END-IF
      *-
      *--* Legge record successivo
           PERFORM LEGG-FILE
           THRU  F-LEGG-FILE
           .
       F-ELAB.
           EXIT.
      *================================================================*
       CTRL-MSG.
      *
      *
           EVALUATE PGPF-MSG-TYPE-ID
               WHEN '0610'
                    PERFORM CTRL-FUNC-CODE
                    THRU  F-CTRL-FUNC-CODE
               WHEN OTHER
                    ADD 1                   TO CTR-CONT-LETTI-SCAR
                    SET WS-SALT-ELAB-SI TO TRUE
           END-EVALUATE
           .
       F-CTRL-MSG.
           EXIT.
      *================================================================*
       CTRL-FUNC-CODE.
      *
DBG==>*    DISPLAY 'CTRL-FUNC-CODE '
DBG==>*    DISPLAY 'PGPF-FUNCT-CODE : ' PGPF-FUNCT-CODE
           EVALUATE PGPF-FUNCT-CODE
               WHEN '697'
                    PERFORM CTRL-FUNC-CODE-697
                    THRU  F-CTRL-FUNC-CODE-697
               WHEN '695'
                    PERFORM CTRL-FUNC-CODE-695
                    THRU  F-CTRL-FUNC-CODE-695
               WHEN '200'
               WHEN '300'
R08421         WHEN '301'
                    PERFORM CTRL-FUNC-CODE-200-300
                    THRU  F-CTRL-FUNC-CODE-200-300
               WHEN OTHER
                    PERFORM CTRL-FCOD-ALTR
                    THRU  F-CTRL-FCOD-ALTR
           END-EVALUATE
           .
       F-CTRL-FUNC-CODE.
           EXIT.
      *================================================================*
       CTRL-FUNC-CODE-697.
      *
           EVALUATE PGPFH-ACTION-CODE
               WHEN '680'
                    PERFORM CTRL-FUNC-CODE-697-A680
                    THRU  F-CTRL-FUNC-CODE-697-A680
               WHEN '681'
                    PERFORM CTRL-FUNC-CODE-697-A681
                    THRU  F-CTRL-FUNC-CODE-697-A681
               WHEN OTHER
                    PERFORM IMPO-ERRO-TEST-ACTI-CODE
                    THRU  F-IMPO-ERRO-TEST-ACTI-CODE
           END-EVALUATE
           .
       F-CTRL-FUNC-CODE-697.
           EXIT.
      *================================================================*
       CTRL-FUNC-CODE-697-A680.
      *
           SET ULTI-RECO-TEST-680           TO TRUE
      *
           ADD 1                            TO CTR-CONT-LETTI-HEAD-680
      *
           IF WS-PRIM-VOLT-SI
              SET WS-PRIM-VOLT-NO           TO TRUE
           END-IF
           .
       F-CTRL-FUNC-CODE-697-A680.
           EXIT.
      *================================================================*
       CTRL-FUNC-CODE-697-A681.
      *
           SET ULTI-RECO-TEST-681           TO TRUE
      *
           ADD 1                            TO CTR-CONT-LETTI-HEAD-681
      *
      *--* Controlla che numero di testa sia maggiore a quello di coda
           IF CTR-CONT-LETTI-TRAIL < CTR-CONT-LETTI-HEAD-681
              CONTINUE
           ELSE
              PERFORM   IMPO-ERRO-TEST
                 THRU F-IMPO-ERRO-TEST
           END-IF
      *
      *--* Controlla il campo data della testata
           PERFORM CTRL-DATA-TEST         THRU F-CTRL-DATA-TEST
      *
           IF WS-PRIM-VOLT-SI
              SET WS-PRIM-VOLT-NO           TO TRUE
           END-IF
           .
       F-CTRL-FUNC-CODE-697-A681.
           EXIT.
      *================================================================*
       CTRL-DATA-TEST.
      *
           MOVE PGPFH-DATE-TIME-CRE(1:8) TO  YYCWUTDA-DATA-CORRENTE(2)
           MOVE 2                        TO  YYCWUTDA-FLAG-SCELTA
      *
           CALL YYUCDATA                 USING YYCWUTDA
      *
           IF YYCWUTDA-FLAG-ERRORE  NOT = ZERO
              PERFORM   IMPO-ERRO-ROUT-TEST
                 THRU F-IMPO-ERRO-ROUT-TEST
           END-IF
           .
       F-CTRL-DATA-TEST.
           EXIT.
      *================================================================*
       CTRL-FUNC-CODE-695.
      *
           SET ULTI-RECO-CODA               TO TRUE
           SET WS-SALT-ELAB-SI              TO TRUE
           ADD 1                            TO CTR-CONT-LETTI-TRAIL
      *
      *--* Controlla che il primo record sia la testata
           IF WS-PRIM-VOLT-SI
              PERFORM   IMPO-ERRO-FCOD-TEST
                 THRU F-IMPO-ERRO-FCOD-TEST
           END-IF
      *
           IF CTR-CONT-LETTI-TRAIL = CTR-CONT-LETTI-HEAD-681
              CONTINUE
           ELSE
              PERFORM   IMPO-ERRO-CODA
                 THRU F-IMPO-ERRO-CODA
           END-IF
           .
       F-CTRL-FUNC-CODE-695.
           EXIT.
      *================================================================*
       CTRL-FUNC-CODE-200-300.
      *
DBG==>*    DISPLAY 'CTRL-FUNC-CODE-200-300 '
           SET ULTI-RECO-DETT               TO TRUE
R14316     SET WS-PAN-PPAY-EVOL-BUSI-NO     TO TRUE
R14316     SET WS-IBAN-ATTI-NO              TO TRUE
R14316*
R14316*    MOVE PGPF-BANK-ACCOUNT(1:27)     TO WS-IBAN
R05818*  inizio parte asteriscata
R14316*    IF  WS-IBAN-ABI = '07601'
R14316*    AND WS-IBAN-CAB = '05138'
R14316*        SET WS-PAN-PPAY-EVOL-BUSI-SI TO TRUE
R14316*    END-IF
R05818*  fine  parte asteriscata
R05818     PERFORM CHIAMA-YPRCFAMI
R05818        THRU EX-CHIAMA-YPRCFAMI
           IF WS-IBAN NOT = PGPF-BANK-ACCOUNT(1:27)
              DISPLAY '********************************************'
              DISPLAY 'ATTENZIONE iban disallineato smac/anagFA    '
              DISPLAY 'IBAN SMAC   = ' PGPF-BANK-ACCOUNT(1:27)
              DISPLAY 'IBAN AnagFA = ' WS-IBAN
              DISPLAY '********************************************'
           END-IF
R05818     IF WS-CC-BANCARIO
R05818        MOVE  PGPF-PAYMT-GEN-DATE  TO WS-PAYMT-GEN-DATE
R05818        IF  PGPF-BUSINESS-DATE NOT = ZERO
R05818            MOVE  PGPF-BUSINESS-DATE   TO WS-BUSINESS-DATE
R05818        END-IF
R05818        MOVE  PGPF-VALUE-DATE      TO WS-VALUE-DATE
R05818        MOVE  PGPF-REFER-PERIOD-DATE-200
R05818                                 TO WS-REFER-PERIOD-DATE-200
R05818     END-IF
      *
R05316     PERFORM   SELE-TABE-GEP-PGP
R05316        THRU F-SELE-TABE-GEP-PGP
      *
R05316     IF YPCRTPGP-DATI-TIPO-ELAB = 'IBAN'
FIANNH     OR YPCRTPGP-DATI-TIPO-ELAB = 'IBAB'
DBG==>*       DISPLAY 'PGPF-FUNCT-CODE ' PGPF-FUNCT-CODE
R05316        IF PGPF-FUNCT-CODE = '200'
R05316           ADD 1             TO CTR-CONT-LETTI-DATI-200
R05316        ELSE
R08421           IF PGPF-FUNCT-CODE = '301'
R08421              ADD 1             TO CTR-CONT-LETTI-DATI-301
R08421           ELSE
R05316              ADD 1             TO CTR-CONT-LETTI-DATI-300
R05316           END-IF
R05316        END-IF
R05316     ELSE
R05316        ADD  1               TO CTR-CONT-LETTI-SCAR
R05316        ADD  1               TO CTR-CONT-LETTI-SCAR-BILL
R05316        ADD  PGPF-PAYMT-TOT  TO CTR-CONT-LETTI-SCAR-BILLI
R00317        SET  WS-TIPO-BILL-SI TO TRUE
R05316        SET  WS-SALT-CONT-SI TO TRUE
R05921*--* In caso di 'BOLL' non deve scrivere la tabella YPTBPGPF
R05921*--* la scriver{ il pgm dedicato ai bolli YPBCEPGB
R05921        IF YPCRTPGP-DATI-TIPO-ELAB = 'BOLL'
R05921           SET  WS-SALT-ELAB-SI TO TRUE
R05921        END-IF
R05316     END-IF.
DBG==>*    DISPLAY 'WS-SALT-ELAB     -3-  ('WS-SALT-ELAB')'

TEST       IF WS-SALT-ELAB-NO
TEST          PERFORM CTRL-ESIS-RIGA      THRU F-CTRL-ESIS-RIGA
TEST       END-IF
      *
DBG==>*    DISPLAY 'WS-SALT-ELAB     -4-  (' WS-SALT-ELAB')'
R08421     IF  WS-SALT-ELAB-NO
R08421     AND PGPF-FUNCT-CODE = '301'
DBG==>*       DISPLAY 'VADO AL CONTROLLO DEL SUMM-ID'
R08421        PERFORM CTRL-SUMM-ID        THRU F-CTRL-SUMM-ID
R08421        GO TO F-CTRL-FUNC-CODE-200-300
R08421     END-IF
      *
           IF  WS-SALT-ELAB-NO
ZANCHI     AND WS-TIPO-BILL-NO
              PERFORM CTRL-IBAN           THRU F-CTRL-IBAN
           END-IF
      *
           IF WS-SALT-ELAB-NO
              PERFORM CTRL-IMPO           THRU F-CTRL-IMPO
           END-IF
R14316*
R14316     IF  WS-PAN-PPAY-EVOL-BUSI-SI
R14316     AND WS-IBAN-ATTI-SI
R14316        PERFORM CTRL-SALDO-E-CUMU   THRU F-CTRL-SALDO-E-CUMU
R14316     END-IF
           .
       F-CTRL-FUNC-CODE-200-300.
           EXIT.
      *================================================================*
       CTRL-ESIS-RIGA.
      *
***********mOVE PGPF-FUNCT-CODE             TO HV-FUNCTION-CODE
           MOVE PGPF-SUMMARY-UID-301        TO HV-SUMMARY-UID
           MOVE PGPF-PAYEMT-UID             TO HV-PAYMENT-UID
DBG==>*    DISPLAY 'CTRL-ESIS-RIGA '
DBG==>*    DISPLAY 'HV-SUMMARY-UID : ' HV-SUMMARY-UID
DBG==>*    DISPLAY 'HV-PAYMENT-UID : ' HV-PAYMENT-UID
      *
           EXEC SQL
                SELECT FUNCTION_CODE
                  INTO :HV-FUNCTION-CODE
                  FROM YPTBPGPF
                 WHERE
                       SUMMARY_UID   = :HV-SUMMARY-UID
                   AND PAYMENT_UID   = :HV-PAYMENT-UID
           END-EXEC
DBG==>*    DISPLAY 'SQLCODE    : ' SQLCODE
      *
           IF SQLCODE = 0
              ADD 1                         TO W300-SCART-DUPKEY
              SET WS-SCRI-SCAR-SI           TO TRUE
              SET WS-SALT-ELAB-SI           TO TRUE
              PERFORM   IMPO-ERRO-X-DUP-KEY
                 THRU F-IMPO-ERRO-X-DUP-KEY
           END-IF.
      *
           .
       F-CTRL-ESIS-RIGA.
           EXIT.
      *================================================================*
       CTRL-IBAN.
      *
           IF PGPF-BANK-ACCOUNT = SPACES
               MOVE SPACES                  TO YPCWS001-RIGA
               STRING 'ACQUIRER ID: ' PGPF-ACQ-ID-CODE
                      ' - IBAN NON VALORIZZATO - RECORD SCARTATO'
                 DELIMITED BY SIZE
                 INTO YPCWS001-RIGA
               END-STRING
               PERFORM SCRIVI-ST          THRU EX-SCRIVI-ST
               SET WS-SCRI-SCAR-SI          TO TRUE
               SET WS-SALT-ELAB-SI          TO TRUE
      *-
      *--* Imposta area x messaggio errori via mail
               MOVE 'Iban non valorizzato su record     -'
                                            TO WS-AREA-APPO-YPOE-DESC
R14316         GO TO F-CTRL-IBAN
           END-IF
R14316*
R14316*--* Se movimento eseguito con carta Postepay Evolution Business
R14316*--* controllo se IBAN � attivo e recupero l'alias della carta
R14316     IF WS-PAN-PPAY-EVOL-BUSI-SI
R14316         PERFORM Z3-CONTROLLA-IBAN-AT
R14316         THRU EX-Z3-CONTROLLA-IBAN-AT
R14316     END-IF
           .
       F-CTRL-IBAN.
           EXIT.
      *==============================================================*
R14316 Z3-CONTROLLA-IBAN-AT.
      *
           INITIALIZE Z3CLSPO2.
           INITIALIZE Z3CWDCOM-DATI-COMUNI.
      *
           MOVE 'IBT'                    TO Z3CWDCOM-FUNZIONE.
           MOVE 'YPBCEPGP'               TO Z3CWDCOM-NOME-PGM.
           MOVE Z3CWDCOM-DATI-COMUNI     TO Z3CLSPO2-DATI-INIZIALI.

           MOVE '0000000'                TO Z3CLSPO2-I-COD-GRUPPO
           MOVE '07601'                  TO Z3CLSPO2-I-COD-ABI-ISTIT
           MOVE '01'                     TO Z3CLSPO2-I-TIPO-CHIAVE.
           MOVE 'IB'                     TO Z3CLSPO2-I-TIPO-SERV.
           MOVE WS-IBAN                  TO Z3CLSPO2-I-ID-CODICE.
DBG==>*    DISPLAY 'WS-IBAN(' WS-IBAN')'

           CALL Z3CWNORB-SPO2-TAB-GEPW   USING Z3CLSPO2

           MOVE Z3CLSPO2-DATI-INIZIALI   TO Z3CWDCOM-DATI-COMUNI.
      *
           EVALUATE Z3CWDCOM-RET-CODE
               WHEN '000'
      *-          IBAN TROVATO
                   SET WS-IBAN-ATTI-SI          TO TRUE

               WHEN '002'
      *-          IBAN NON TROVATO
                   MOVE SPACES                  TO YPCWS001-RIGA
                   STRING ' - IBAN NON TROVATO SU CARD:' WS-IBAN
                     DELIMITED BY SIZE
                     INTO YPCWS001-RIGA
                   END-STRING
                   PERFORM SCRIVI-ST          THRU EX-SCRIVI-ST
                   SET WS-SCRI-SCAR-SI          TO TRUE
                   SET WS-SALT-ELAB-SI          TO TRUE
      *-
      *--* Imposta area x messaggio errori via mail
                   MOVE 'Iban non presente su archivi CARD  -'
                                            TO WS-AREA-APPO-YPOE-DESC

              WHEN OTHER
                   MOVE SPACES                  TO YPCWS001-RIGA
                   STRING ' - ERRORE SU CARD X IBAN:' WS-IBAN
                     DELIMITED BY SIZE
                     INTO YPCWS001-RIGA
                   END-STRING
                   PERFORM SCRIVI-ST          THRU EX-SCRIVI-ST
                   SET WS-SCRI-SCAR-SI          TO TRUE
                   SET WS-SALT-ELAB-SI          TO TRUE
      *-
      *--* Imposta area x messaggio errori via mail
                   MOVE 'Errore su CARD x Iban              -'
                                            TO WS-AREA-APPO-YPOE-DESC

           END-EVALUATE.
      *
       EX-Z3-CONTROLLA-IBAN-AT.
           EXIT.
      *==============================================================*
R14316 Z3-CHIAMA-UIFA.
piero *
           INITIALIZE Z3CLUIFA.
      *
           MOVE '001'                       TO Z3CLUIFA-IN-FUNZ
           MOVE Z3CLSPO2-O-ID-VALORE-CHIAVE TO Z3CLUIFA-IN-ID-REALE

           CALL Z3BCUIFA                 USING Z3CLUIFA.

      *
           EVALUATE Z3CLUIFA-OU-RET-CODE
               WHEN '000'
R13419****  '003'  carta in blocco per CA
R13419         WHEN '003'
                   CONTINUE
R15420****  '004'  carta in blocco per C0 - uso scorretto
R15420         WHEN '004'
R15420           PERFORM CHIAMA-Z3BCUI99
R15420              THRU F-CHIAMA-Z3BCUI99
               WHEN '001'
                   MOVE SPACES                  TO YPCWS001-RIGA
                   STRING ' - PARAMETRI ERRATI SU CARD-'
                          Z3CLUIFA-OU-DESCR-ERR
                          '-Rifer.errore('
                          Z3CLUIFA-OU-RIFE-ERR')'
                     DELIMITED BY SIZE
                     INTO YPCWS001-RIGA
                   END-STRING
                   PERFORM SCRIVI-ST          THRU EX-SCRIVI-ST
                   SET WS-SCRI-SCAR-SI          TO TRUE
                   SET WS-SALT-ELAB-SI          TO TRUE
      *-
      *--* Imposta area x messaggio errori via mail
                   MOVE 'Parametri errati passati a Z3BCUIFA-'
                                            TO WS-AREA-APPO-YPOE-DESC

               WHEN '002'
                   MOVE SPACES                  TO YPCWS001-RIGA
                   STRING ' - Carta non attiva - '
                          ' Ret.code('
                          Z3CLUIFA-OU-RET-CODE
                          ')'
                          '- Rif.err.('
                          Z3CLUIFA-OU-RIFE-ERR')'
                     DELIMITED BY SIZE
                     INTO YPCWS001-RIGA
                   END-STRING
                   PERFORM SCRIVI-ST          THRU EX-SCRIVI-ST
                   SET WS-SCRI-SCAR-SI          TO TRUE
                   SET WS-SALT-ELAB-SI          TO TRUE
      *-
      *--* Imposta area x messaggio errori via mail
                   MOVE 'Carta non attiva                   -'
                                            TO WS-AREA-APPO-YPOE-DESC

              WHEN OTHER
                   MOVE SPACES                  TO YPCWS001-RIGA
                   STRING ' - ELABORAZIONE KO-'
                          Z3CLUIFA-OU-DESCR-ERR
                          '-Rifer.errore('
                          Z3CLUIFA-OU-RIFE-ERR')'
                   DELIMITED BY SIZE
                   INTO YPCWS001-RIGA
                   END-STRING
                   PERFORM SCRIVI-ST          THRU EX-SCRIVI-ST
                   SET WS-SCRI-SCAR-SI          TO TRUE
                   SET WS-SALT-ELAB-SI          TO TRUE
      *-
      *--* Imposta area x messaggio errori via mail
                   MOVE 'Errore sul programma -Z3BCUIFA     -'
                                            TO WS-AREA-APPO-YPOE-DESC

           END-EVALUATE.
      *
       EX-Z3-CHIAMA-UIFA.
           EXIT.
      *================================================================*
       CTRL-IMPO.
      *
            IF PGPF-PAYMT-TOT = 0
               MOVE SPACES                  TO YPCWS001-RIGA
                STRING  'IBAN: ' PGPF-BANK-ACCOUNT
                        ' - IMPORTO NON VALORIZZATO'
                        DELIMITED BY SIZE
                        INTO YPCWS001-RIGA
                END-STRING
                PERFORM SCRIVI-ST         THRU EX-SCRIVI-ST
                SET WS-SCRI-SCAR-SI         TO TRUE
                SET WS-SALT-ELAB-SI         TO TRUE
R03817          SET WS-IMPO-ZERO-SI         TO TRUE
      *-
      *--* Imposta area x messaggio errori via mail
                MOVE 'Importo non valorizzato su record  -'
                                            TO WS-AREA-APPO-YPOE-DESC
            END-IF
           .
       F-CTRL-IMPO.
           EXIT.
      *================================================================*
R08421 CTRL-SUMM-ID.
R08421*
R08421      IF PGPF-SUMMARY-UID-301 = 0
DBG==>*        DISPLAY 'SUMMARY ID ZERO'
R08421         MOVE SPACES                  TO YPCWS001-RIGA
R08421          STRING  'CAMPO SUMMARY ID: ' PGPF-SUMMARY-UID-301
R08421                  ' - MPORTO NON VALORIZZATO'
R08421                  DELIMITED BY SIZE
R08421                  INTO YPCWS001-RIGA
R08421          END-STRING
R08421          PERFORM SCRIVI-ST         THRU EX-SCRIVI-ST
R08421          SET WS-SCRI-SCAR-SI         TO TRUE
R08421          SET WS-SALT-ELAB-SI         TO TRUE
R08421*-
R08421*--* Imposta area x messaggio errori via mail
R08421          MOVE 'Campo summary id non valorizzato su record  -'
R08421                                      TO WS-AREA-APPO-YPOE-DESC
R08421      END-IF
R08421     .
R08421 F-CTRL-SUMM-ID.
           EXIT.
      *
      *================================================================*
R14316 CTRL-SALDO-E-CUMU.
R14316*
R14316*--* Recupera pan di II, saldo e capacit� nominale della carta
R14316     PERFORM Z3-CHIAMA-UIFA
R14316     THRU EX-Z3-CHIAMA-UIFA
      *-
      *--* Se Importo movimento dare superiore a saldo monte moneta
      *--* il programma scarta il record creando DCD di apertura
Zanchi*
Zanchi*    IF  WS-SALT-ELAB-NO
Zanchi*    AND PGPF-DB-CR-FLAG = 'D'
Zanchi*        COMPUTE WS-IMPORTO-MOV = PGPF-PAYMT-TOT / 100
Zanchi*        IF WS-IMPORTO-MOV > Z3CLUIFA-OU-SALDO-DISP
Zanchi*           MOVE SPACES                 TO YPCWS001-RIGA
Zanchi*           MOVE WS-IMPORTO-MOV         TO WS-IMPORTO-MOV-D
Zanchi*           MOVE Z3CLUIFA-OU-SALDO-DISP TO WS-SALDO-D
Zanchi*             STRING ' - IMPORTO DARE('WS-IMPORTO-MOV-D')'
Zanchi*             ' MAGGIORE DEL SALDO('WS-SALDO-D')'
Zanchi*             DELIMITED BY SIZE
Zanchi*             INTO YPCWS001-RIGA
Zanchi*           END-STRING
Zanchi*           PERFORM SCRIVI-ST        THRU EX-SCRIVI-ST
Zanchi*           SET WS-SCRI-SCAR-SI        TO TRUE
Zanchi*           SET WS-SALT-ELAB-SI        TO TRUE
Zanchi*-
Zanchi*--* Imposta area x messaggio errori via mail
Zanchi*           MOVE 'Importo dare maggiore del saldo    -'
Zanchi*                                      TO WS-AREA-APPO-YPOE-DESC
Zanchi*        END-IF
Zanchi*    END-IF
      *-
      *--* Se Importo movimento avere + saldo monte moneta � superiore
      *--* alla capacit� nominale della carta
      *--* il programma scarta il record creando DCD di apertura
           IF  WS-SALT-ELAB-NO
           AND PGPF-DB-CR-FLAG = 'C'
               COMPUTE WS-IMPORTO-MOV  =  PGPF-PAYMT-TOT / 100
               COMPUTE WS-IMPORTO-TOT  =  WS-IMPORTO-MOV
                                       +  Z3CLUIFA-OU-SALDO-DISP
               IF WS-IMPORTO-TOT > Z3CLUIFA-OU-CAP-NOMIN
                  MOVE SPACES                TO YPCWS001-RIGA
                  STRING ' - IMPORTO AVERE + SALDO '
                         ' MAGGIORE DELLA CAPACITA''NOMINALE CARTA'
                    DELIMITED BY SIZE
                    INTO YPCWS001-RIGA
                  END-STRING
                  PERFORM SCRIVI-ST        THRU EX-SCRIVI-ST
R07420            SET WS-CONTAB-SI           TO TRUE
                  SET WS-SCRI-SCAR-SI        TO TRUE
                  SET WS-SALT-ELAB-SI        TO TRUE
      *-
      *--* Imposta area x messaggio errori via mail
                  MOVE 'Imp.avere + saldo > capac.nom.carta-'
                                             TO WS-AREA-APPO-YPOE-DESC
               END-IF
           END-IF
           .
       F-CTRL-SALDO-E-CUMU.
           EXIT.
      *
      *================================================================*
       CTRL-FCOD-ALTR.
      *-
      *--* Controlla che il primo record sia la testata
           IF WS-PRIM-VOLT-SI
              PERFORM   IMPO-ERRO-FCOD-TEST
                 THRU F-IMPO-ERRO-FCOD-TEST
           ELSE
              IF PGPFH-FUNCT-CODE IS NOT NUMERIC
                 PERFORM   IMPO-ERRO-FCOD
                    THRU F-IMPO-ERRO-FCOD
              END-IF
           END-IF
      *
           ADD 1                            TO CTR-CONT-LETTI-SCAR
           .
       F-CTRL-FCOD-ALTR.
           EXIT.
      *================================================================*
       ACQU-DATI.
      *
           EVALUATE PGPF-FUNCT-CODE
               WHEN '697'
                    PERFORM ACQU-DATI-697
                    THRU  F-ACQU-DATI-697
               WHEN '200'
               WHEN '300'
                    PERFORM ACQU-DATI-200-300
                    THRU  F-ACQU-DATI-200-300
               WHEN '301'
                    CONTINUE
               WHEN OTHER
                    SET WS-SALT-ELAB-SI     TO TRUE
           END-EVALUATE
           .
       F-ACQU-DATI.
           EXIT.
      *================================================================*
       ACQU-DATI-697.
      *
           SET WS-SALT-ELAB-SI              TO TRUE
      *
           EVALUATE PGPFH-ACTION-CODE
               WHEN '681'
                    MOVE PGPFH-DATE-TIME-CRE TO COM-DATE-TIME-N
               WHEN OTHER
                    CONTINUE
           END-EVALUATE
           .
       F-ACQU-DATI-697.
           EXIT.
      *================================================================*
       ACQU-DATI-200-300.
      *
R05818*R14316     IF WS-PAN-PPAY-EVOL-BUSI-NO
R05818     IF WS-CC-POSTE
           AND WS-TIPO-BILL-NO
TEST          PERFORM RICERCA-FRAZIONARIO
TEST          THRU EX-RICERCA-FRAZIONARIO
TEST  *       MOVE '12345'            TO COMSD50-FILIALE
TEST  *       MOVE 123456789012       TO COMSD50-RAPPORT
TESt  *       MOVE '1234'             TO COMSD50-CATRAPP
TEST  *       MOVE '55355'            TO COMSD50-FILIALE
TEST  *       MOVE 191003151737       TO COMSD50-RAPPORT
TESt  *       MOVE '1234'             TO COMSD50-CATRAPP
R14316     END-IF
      *
R00317*--* Richiesta di Zanchi con ticket numero 2017011210000069
R00317*--* Se tipo flusso Billing (invio flusso fatturazione)
R00317*--* non devono essere impostate le causali e il cod.operazione e
R00317*--* il flusso deve essere scritto in tabella PGPF e non scartati
R00317     IF WS-TIPO-BILL-SI
R00317        MOVE SPACES           TO W100-CAUSALE
R00317                                 W100-CODOPE
R00317     ELSE
              IF WS-SCRI-SCAR-NO
R12117*          PERFORM RICERCA-CAUS-OPE
R12117*          THRU EX-RICERCA-CAUS-OPE
R12117           PERFORM RECUPERA-FLAG-TIPO-POS
R12117           THRU EX-RECUPERA-FLAG-TIPO-POS
R12117           IF  WS-SCRI-SCAR-NO
R12117              PERFORM RICERCA-CAUS-OPE-NEW
R12117              THRU EX-RICERCA-CAUS-OPE-NEW
R12117           END-IF
              END-IF
R00317     END-IF
R14316*
R14316     IF  WS-SCRI-SCAR-NO
R14316     AND WS-PAN-PPAY-EVOL-BUSI-SI
ZANCHI     AND WS-TIPO-BILL-NO
R14316        PERFORM CERCA-PAN-3TRA
R14316        THRU  F-CERCA-PAN-3TRA
R14316     END-IF
           .
       F-ACQU-DATI-200-300.
           EXIT.
      *================================================================*
       RICERCA-FRAZIONARIO.
      *
           MOVE SPACES                      TO COMSD50-FILIALE
                                               COMSD50-CATRAPP
           MOVE 0                           TO COMSD50-RAPPORT
           INITIALIZE INCC-CV20
      *
           MOVE WS-IBAN(16:12)              TO INCC-CV20-RAPPORT
                                               WS-COM-RAPPORTO
DBG==>*    display 'ricerca-frazionario '
DBG==>*    display 'INCC-CV20-RAPPORTO: ' INCC-CV20-RAPPORT
           CALL CRVYD228                 USING PPTCINCC
      *
DBG==>*    display 'INCC-RETCODE: ' INCC-RETCODE
           EVALUATE INCC-RETCODE
               WHEN HIGH-VALUES
                    PERFORM IMPO-ERRO-INCC-HIVA
                    THRU  F-IMPO-ERRO-INCC-HIVA
               WHEN 'NO'
                    PERFORM IMPO-ERRO-INCC-NO
                    THRU  F-IMPO-ERRO-INCC-NO
               WHEN 'NF'
                    PERFORM IMPO-ERRO-INCC-NF
                    THRU  F-IMPO-ERRO-INCC-NF
                    SET WS-SCRI-SCAR-SI     TO TRUE
                    SET WS-SALT-ELAB-SI     TO TRUE
      *-
      *--* Imposta area x messaggio errori via mail
                    MOVE 'Notfound da routine CRVYD228  -'
                                            TO WS-AREA-APPO-YPOE-DESC
               WHEN OTHER
                    MOVE INCC-CV20-FILIALE  TO COMSD50-FILIALE
                    MOVE INCC-CV20-RAPPORT  TO COMSD50-RAPPORT
                    MOVE INCC-CV20-CATRAPP  TO COMSD50-CATRAPP
           END-EVALUATE
           .
       EX-RICERCA-FRAZIONARIO.
           EXIT.
      *================================================================*
R05818 RICERCA-FRAZIONARIO-CC.
      *
           MOVE SPACES                      TO COMSD50-FILIALE
                                               COMSD50-CATRAPP
           MOVE 0                           TO COMSD50-RAPPORT
           INITIALIZE INCC-CV20
      *
           MOVE YPCRTCCB-RAPPORTO           TO INCC-CV20-RAPPORT
                                               WS-COM-RAPPORTO
DBG==>*    display 'ricerca-frazionario-cc'
DBG==>*    display 'INCC-CV20-RAPPORTO: ' INCC-CV20-RAPPORT
           CALL CRVYD228                 USING PPTCINCC
      *
DBG==>*    display 'INCC-RETCODE :' INCC-RETCODE
           EVALUATE INCC-RETCODE
               WHEN HIGH-VALUES
                    PERFORM IMPO-ERRO-INCC-HIVA
                    THRU  F-IMPO-ERRO-INCC-HIVA
               WHEN 'NO'
               WHEN 'NF'
                    PERFORM IMPO-ERRO-INCC-NO
                    THRU  F-IMPO-ERRO-INCC-NO
               WHEN OTHER
                    MOVE INCC-CV20-FILIALE  TO COMSD50-FILIALE
                    MOVE INCC-CV20-RAPPORT  TO COMSD50-RAPPORT
                    MOVE INCC-CV20-CATRAPP  TO COMSD50-CATRAPP
           END-EVALUATE
           .
R05818 EX-RICERCA-FRAZIONARIO-CC.
           EXIT.
      *================================================================*
       RICERCA-CAUS-OPE.
      *
            MOVE SPACES                     TO W100-CAUSALE W100-CODOPE
            MOVE PGPF-BANK-ACC-TYP          TO W100-BANK-ACC-TYPE
            MOVE PGPF-PAYMT-TYPE            TO W100-BANK-ACC-TYPE(8:3)
      *
R14316      IF WS-PAN-PPAY-EVOL-BUSI-NO
               PERFORM VARYING IND1 FROM 1 BY 1
                       UNTIL   IND1 > VALORE-MAX
                IF W100-BANK-ACC-TYPE = YPCRTPI-OLI-AE-PROVENIENZA(IND1)
                   IF PGPF-DB-CR-FLAG = 'D'
                      MOVE YPCRTPI-OLI-AE-CAUSALE-ADD(IND1)
                                               TO  W100-CAUSALE
                      MOVE YPCRTPI-OLI-AE-CODOPE-ADD(IND1)
                                               TO  W100-CODOPE
                    ELSE
                      MOVE YPCRTPI-OLI-AE-CAUSALE-ACC(IND1)
                                               TO  W100-CAUSALE
                      MOVE YPCRTPI-OLI-AE-CODOPE-ACC(IND1)
                                              TO  W100-CODOPE
                   END-IF
                   ADD VALORE-MAX TO IND1
                END-IF
               END-PERFORM
R14316      ELSE
      *
R14316         PERFORM VARYING IND1 FROM 1 BY 1
R14316                 UNTIL   IND1 > VALORE-MAX
R14316          IF W100-BANK-ACC-TYPE = WSCRTPI-OLI-AE-PROVENIENZA(IND1)
R14316             IF PGPF-DB-CR-FLAG = 'D'
R14316                MOVE WSCRTPI-OLI-AE-CAUSALE-ADD(IND1)
R14316                                         TO  W100-CAUSALE
R14316                MOVE WSCRTPI-OLI-AE-CODOPE-ADD(IND1)
R14316                                         TO  W100-CODOPE
R14316              ELSE
R14316                MOVE WSCRTPI-OLI-AE-CAUSALE-ACC(IND1)
R14316                                         TO  W100-CAUSALE
R14316                MOVE WSCRTPI-OLI-AE-CODOPE-ACC(IND1)
R14316                                        TO  W100-CODOPE
R14316             END-IF
R14316             ADD VALORE-MAX TO IND1
R14316          END-IF
R14316         END-PERFORM
R14316      END-IF
R14316      .
      *
            IF  W100-CAUSALE = SPACES
                MOVE SPACES                 TO YPCWS001-RIGA
                STRING  'IBAN: ' PGPF-BANK-ACCOUNT
                        ' CAUSALE NON TROVATA SU YYDTABE'
                        DELIMITED BY SIZE
                        INTO YPCWS001-RIGA
                END-STRING
                PERFORM SCRIVI-ST         THRU EX-SCRIVI-ST
                SET WS-SCRI-SCAR-SI         TO TRUE
                SET WS-SALT-ELAB-SI         TO TRUE
      *-
      *--* Imposta area x messaggio errori via mail
                MOVE 'Causale non trovata                -'
                                            TO WS-AREA-APPO-YPOE-DESC
            END-IF
      *
            IF  W100-CODOPE  = SPACES
                MOVE SPACES                 TO YPCWS001-RIGA
                STRING  'IBAN: ' PGPF-BANK-ACCOUNT
                        ' CODOPE  NON TROVATO SU YYDTABE'
                        DELIMITED BY SIZE
                        INTO YPCWS001-RIGA
                END-STRING
                PERFORM SCRIVI-ST         THRU EX-SCRIVI-ST
                SET WS-SCRI-SCAR-SI         TO TRUE
                SET WS-SALT-ELAB-SI         TO TRUE
      *-
      *--* Imposta area x messaggio errori via mail
                MOVE 'Codice operatore non trovato       -'
                                            TO WS-AREA-APPO-YPOE-DESC
             ELSE
                MOVE W100-CODOPE            TO CRVSD50-CODOPE
            END-IF
            .
       EX-RICERCA-CAUS-OPE.
            EXIT.
      *================================================================*
R12117 RECUPERA-FLAG-TIPO-POS.
      *
R05818*     MOVE SPACE  TO  WK-FLAG-TIPO-POS
            IF WS-PAN-PPAY-EVOL-BUSI-NO
             IF PGPF-LEVEL-PAY-CODE = 'S' OR 'C'
R05818*         PERFORM CHIAMA-YPRCFAMI
R05818*            THRU EX-CHIAMA-YPRCFAMI
R05818             CONTINUE
             ELSE
              MOVE SPACES                 TO YPCWS001-RIGA
              STRING  'PGPF-LEVEL-PAY-CODE: ' PGPF-LEVEL-PAY-CODE
                        ' NON PREVISTA '
                        DELIMITED BY SIZE
                        INTO YPCWS001-RIGA
               END-STRING
               PERFORM SCRIVI-ST         THRU EX-SCRIVI-ST
               SET WS-SCRI-SCAR-SI         TO TRUE
               SET WS-SALT-ELAB-SI         TO TRUE
      *-
      *--* Imposta area x messaggio errori via mail
               MOVE 'LEVEL-PAY-CODE non prevista   -'
                                            TO WS-AREA-APPO-YPOE-DESC
             END-IF
            ELSE
              MOVE 'F' TO WK-FLAG-TIPO-POS
            END-IF
            .
R12117 EX-RECUPERA-FLAG-TIPO-POS.
            EXIT.
      *================================================================*
R12117 CHIAMA-YPRCFAMI.
      *
R05818      MOVE SPACE  TO  WK-FLAG-TIPO-POS
R05818      SET WS-CC-POSTE              TO TRUE
            INITIALIZE YPCWFAMI
            .
            IF PGPF-LEVEL-PAY-CODE = 'C'
               MOVE '01'                  TO  YPCWFAMI-I-FUNZ
            ELSE
               MOVE '02'                  TO  YPCWFAMI-I-FUNZ
            END-IF
            .
            MOVE PGPF-ME-ID-CODE-200      TO  YPCWFAMI-I-DATO
DBG==>*     DISPLAY 'PGPF-ME-ID-CODE-200 ('PGPF-ME-ID-CODE-200')'
DBG==>*     DISPLAY 'YPCWFAMI-I-FUNZ     ('YPCWFAMI-I-FUNZ  ')'
DBG==>*     DISPLAY 'YPCWFAMI-I-DATO     ('YPCWFAMI-I-DATO  ')'
            CALL  YPRCFAMI        USING YPCWFAMI
DBG==>*     DISPLAY 'DOPO CALL YPRCFAMI  ('YPCWFAMI-O-ESIT ')'
DBG==>*     DISPLAY 'YPCWFAMI-O-TIPO-SP-A('YPCWFAMI-O-TIPO-SP-ACCR')'
            IF YPCWFAMI-O-ESIT NOT = 'OK'
171122          DISPLAY 'YPCWFAMI-O-ESIT    ('YPCWFAMI-O-ESIT  ')'
171122          DISPLAY 'YPCWFAMI-I-FUNZ    ('YPCWFAMI-I-FUNZ  ')'
171122          DISPLAY 'YPCWFAMI-I-DATO    ('YPCWFAMI-I-DATO  ')'
171122*--* Per gestire i primi esercenti FA non censiti correttamente
171122*--* in anagrafica di forza F
171122          MOVE 'F'       TO WK-FLAG-TIPO-POS
                MOVE PGPF-BANK-ACCOUNT(1:27) TO WS-IBAN
R07419*--* Nel caso il tipo conto non corrisponde ad un valore atteso,
R07419*--* imposta le variabili facendo logica sull'IBAN d'input PGPF
R07419          PERFORM   GEST-IBAN-CASO-ANOM
R07419             THRU F-GEST-IBAN-CASO-ANOM
            ELSE
                  IF YPCWFAMI-O-PROD-E-PROD(1) = 'VIRTUAL POS'
                     MOVE 'V'       TO WK-FLAG-TIPO-POS
                  ELSE
                     MOVE 'F'       TO WK-FLAG-TIPO-POS
                  END-IF

                  IF PGPF-DB-CR-FLAG = 'C'

101218               IF YPCWFAMI-O-IBAN-ACCR NOT = SPACES
                        MOVE YPCWFAMI-O-IBAN-ACCR TO WS-IBAN
101218               ELSE
101218                  MOVE PGPF-BANK-ACCOUNT(1:27) TO WS-IBAN
101218               END-IF

R05818               EVALUATE  YPCWFAMI-O-TIPO-SP-ACCR
R05818                 WHEN '02'
R05818                    SET WS-CC-POSTE            TO TRUE
R05818                 WHEN '01'
R05818                    SET WS-PP-EVO                TO TRUE
R05818                    SET WS-PAN-PPAY-EVOL-BUSI-SI TO TRUE
R05818                 WHEN '03'
R05818                    SET WS-CC-BANCARIO        TO TRUE
R07419                 WHEN OTHER
R07419*--* Nel caso il tipo conto non corrisponde ad un valore atteso,
R07419*--* imposta le variabili facendo logica sull'IBAN d'input PGPF
R07419                    PERFORM   GEST-IBAN-CASO-ANOM
R07419                       THRU F-GEST-IBAN-CASO-ANOM
R05818               END-EVALUATE
                  ELSE

101218               IF YPCWFAMI-O-IBAN-ADDE NOT = SPACES
                        MOVE YPCWFAMI-O-IBAN-ADDE TO WS-IBAN
101218               ELSE
101218                  MOVE PGPF-BANK-ACCOUNT(1:27) TO WS-IBAN
101218               END-IF

                     EVALUATE  YPCWFAMI-O-TIPO-SP-ADDE
                       WHEN '02'
                          SET WS-CC-POSTE            TO TRUE
                       WHEN '01'
                          SET WS-PP-EVO                TO TRUE
                          SET WS-PAN-PPAY-EVOL-BUSI-SI TO TRUE
                       WHEN '03'
                          SET WS-CC-BANCARIO        TO TRUE
R07419                 WHEN OTHER
R07419*--* Nel caso il tipo conto non corrisponde ad un valore atteso,
R07419*--* imposta le variabili facendo logica sull'IBAN d'input PGPF
R07419                    PERFORM   GEST-IBAN-CASO-ANOM
R07419                       THRU F-GEST-IBAN-CASO-ANOM
                     END-EVALUATE
                  END-IF
            END-IF
            .
R12117 EX-CHIAMA-YPRCFAMI.
            EXIT.
      *================================================================*
R07419 GEST-IBAN-CASO-ANOM.
      *
      *--* Nel caso di Postepay ( ABI/CAB uguali a 07601/05138
      *--* o 36081/05138) imposta conto a livello di PostePay
            IF PGPF-BANK-ACCOUNT(11:05) = '05138'
            AND (   PGPF-BANK-ACCOUNT(06:05) = '07601'
                 OR PGPF-BANK-ACCOUNT(06:05) = '36081'
                 )
               SET WS-PP-EVO                TO TRUE
               SET WS-PAN-PPAY-EVOL-BUSI-SI TO TRUE
            ELSE
      *--* Nel caso di conto BancoPoste ( ABI uguale a 07601)
      *--* imposta conto a livello di Banco Poste
               IF PGPF-BANK-ACCOUNT(06:05) = '07601'
                  SET WS-CC-POSTE           TO TRUE
               ELSE
      *--* Altrimenti imposta conto a livello di Conto Bancario
                  SET WS-CC-BANCARIO        TO TRUE
               END-IF
            END-IF
            .
R07419 F-GEST-IBAN-CASO-ANOM.
R07419     EXIT.
      *================================================================*
R12117 RICERCA-CAUS-OPE-NEW.
      *
            MOVE SPACES                     TO W100-CAUSALE W100-CODOPE
      *
            PERFORM VARYING IND1 FROM 1 BY 1
R12019*                UNTIL   IND1 > 200
R12019                 UNTIL   IND1 > 400
                IF  PGPF-BANK-ACC-TYP = WK-NW-PGPF-BANK-ACC-TYP(IND1)
                AND PGPF-PAYMT-TYPE   = WK-NW-PGPF-PAYMT-TYPE(IND1)
                AND WK-FLAG-TIPO-POS  = WK-NW-FLAG-TIPO-POS(IND1)
                 IF WS-PAN-PPAY-EVOL-BUSI-NO
                   IF PGPF-DB-CR-FLAG = 'D'
                      MOVE WK-NW-CAUSALE-ADD(IND1)
                                               TO  W100-CAUSALE
                      MOVE WK-NW-CAUSALE-COM(IND1)
                                               TO  W100-CODOPE
                    ELSE
                      MOVE WK-NW-CAUSALE-ACC(IND1)
                                               TO  W100-CAUSALE
                      MOVE WK-NW-CAUSALE-COM-ACC(IND1)
                                              TO  W100-CODOPE
                   END-IF
                 ELSE
                   IF PGPF-DB-CR-FLAG = 'D'
                      MOVE WK-NW-CAUSALE-ADD-P(IND1)
                                               TO  W100-CAUSALE
                      MOVE WK-NW-CAUSALE-COM-P(IND1)
                                               TO  W100-CODOPE
                    ELSE
                      MOVE WK-NW-CAUSALE-ACC-P(IND1)
                                               TO  W100-CAUSALE
                      MOVE WK-NW-CAUSALE-COM-ACC-P(IND1)
                                              TO  W100-CODOPE
                   END-IF
                 END-IF
R12019*          ADD 200        TO IND1
R12019           ADD 400        TO IND1
                END-IF
            END-PERFORM
            .
      *
            IF  W100-CAUSALE = SPACES
                MOVE SPACES                 TO YPCWS001-RIGA
                STRING  'IBAN: ' PGPF-BANK-ACCOUNT
                        ' CAUSALE NON TROVATA SU YYDTABE'
                        DELIMITED BY SIZE
                        INTO YPCWS001-RIGA
                END-STRING
                PERFORM SCRIVI-ST         THRU EX-SCRIVI-ST
                SET WS-SCRI-SCAR-SI         TO TRUE
                SET WS-SALT-ELAB-SI         TO TRUE
      *-
      *--* Imposta area x messaggio errori via mail
                MOVE 'Causale non trovata                -'
                                            TO WS-AREA-APPO-YPOE-DESC
171122          DISPLAY 'W100-CAUSALE = SPACES'
171122          DISPLAY 'PGPF-BANK-ACC-TYP  ('PGPF-BANK-ACC-TYP')'
171122          DISPLAY 'PGPF-PAYMT-TYPE    ('PGPF-PAYMT-TYPE  ')'
171122          DISPLAY 'WK-FLAG-TIPO-POS   ('WK-FLAG-TIPO-POS ')'
            END-IF
      *
            IF  W100-CODOPE  = SPACES
                MOVE SPACES                 TO YPCWS001-RIGA
                STRING  'IBAN: ' PGPF-BANK-ACCOUNT
                        ' CODOPE  NON TROVATO SU YYDTABE'
                        DELIMITED BY SIZE
                        INTO YPCWS001-RIGA
                END-STRING
                PERFORM SCRIVI-ST         THRU EX-SCRIVI-ST
                SET WS-SCRI-SCAR-SI         TO TRUE
                SET WS-SALT-ELAB-SI         TO TRUE
      *-
      *--* Imposta area x messaggio errori via mail
                MOVE 'Codice operatore non trovato       -'
                                            TO WS-AREA-APPO-YPOE-DESC
171122          DISPLAY 'W100-CODOPE  = SPACES'
171122          DISPLAY 'PGPF-BANK-ACC-TYP  ('PGPF-BANK-ACC-TYP')'
171122          DISPLAY 'PGPF-PAYMT-TYPE    ('PGPF-PAYMT-TYPE  ')'
171122          DISPLAY 'WK-FLAG-TIPO-POS   ('WK-FLAG-TIPO-POS ')'
             ELSE
                MOVE W100-CODOPE            TO CRVSD50-CODOPE
            END-IF
            .
R12117 EX-RICERCA-CAUS-OPE-NEW.
            EXIT.
      *================================================================*
      * PARTENDO DAL PAN II TRACCIA CERCO IL PAN BANCOMAT
      *================================================================
R14316 CERCA-PAN-3TRA.
      *
           INITIALIZE Z3CLGE90.
      *
           MOVE Z3CLUIFA-OU-PAN-II-TR       TO WS-PAN-EUC
           .
           IF WS-BIN-EUC = '588602'
           OR WS-BIN-EUC = '677035'
              MOVE WS-CARTA-EUC(2:8)        TO Z3CLGE90-PAN
           ELSE
              MOVE WS-PAN-EUC               TO Z3CLGE90-PAN
           .
           MOVE WS-PAN-EUC                  TO Z3CLGE90-ID-GENERICO
           MOVE 'CLE'                       TO Z3CLGE90-FUNZ
           .
           CALL Z3BCGE90 USING Z3CLGE90
           .
      *-
      *--* Se routine code OK va a fine label
           IF Z3CLGE90-OK
              GO TO F-CERCA-PAN-3TRA
           END-IF
           .
      *-
      *--* Se routine code KO diversifica le segnalazioni di errore
           EVALUATE TRUE
              WHEN   Z3CLGE90-PAN-SPACE
                   MOVE SPACES                 TO YPCWS001-RIGA
                   STRING 'ERR. ROUTINE Z3BCGE90 '
                      'PAN II TRACCIA NON VALORIZZATO'
                        DELIMITED BY SIZE
                        INTO YPCWS001-RIGA
                   END-STRING
                   PERFORM SCRIVI-ST         THRU EX-SCRIVI-ST
                   SET WS-SCRI-SCAR-SI         TO TRUE
                   SET WS-SALT-ELAB-SI         TO TRUE
      *-
      *--* Imposta area x messaggio errori via mail
                   MOVE 'Pan II non valorizz- Errore GE90   -'
                                            TO WS-AREA-APPO-YPOE-DESC
               WHEN    Z3CLGE90-PAN-NON-TROVATO
                   MOVE SPACES                 TO YPCWS001-RIGA
                   STRING 'ERR. ROUTINE Z3BCGE90 '
                      'PAN III TRACCIA NON TROVATO'
                        DELIMITED BY SIZE
                        INTO YPCWS001-RIGA
                   END-STRING
                   PERFORM SCRIVI-ST         THRU EX-SCRIVI-ST
                   SET WS-SCRI-SCAR-SI         TO TRUE
                   SET WS-SALT-ELAB-SI         TO TRUE
      *-
      *--* Imposta area x messaggio errori via mail
                   MOVE 'Pan III non trovato- Errore GE90   -'
                                            TO WS-AREA-APPO-YPOE-DESC
               WHEN OTHER
                   MOVE SPACES                 TO YPCWS001-RIGA
                   STRING 'ERR. ROUTINE Z3BCGE90 '
                      'SQLCODE ' Z3CLGE90-SQLCODE
                        DELIMITED BY SIZE
                        INTO YPCWS001-RIGA
                   END-STRING
                   PERFORM SCRIVI-ST         THRU EX-SCRIVI-ST
                   SET WS-SCRI-SCAR-SI         TO TRUE
                   SET WS-SALT-ELAB-SI         TO TRUE
      *-
      *--* Imposta area x messaggio errori via mail
                   MOVE 'Errore generico routine Z3BCGE90   -'
                                            TO WS-AREA-APPO-YPOE-DESC
           END-EVALUATE

           .
       F-CERCA-PAN-3TRA.
           EXIT.
      *================================================================*
       SCRIVI-SCARTI.
      *
            WRITE OUSCARTI-REC            FROM YPCRPGPF-WORK
      *
            IF NOT OUSC-NORMAL
               MOVE ST-OUSCARTI             TO P303-FILE-STATUS
               MOVE '08'                    TO P303-MSGER-RIF
               MOVE 'OUSCARTI'              TO P303-MSGER-FILE
               MOVE 'WRITE'                 TO P303-MSGER-TIPO
               MOVE 'ERRORE WRITE FILE ANOM' TO P303-MSGER-DESCR
               PERFORM ERRORE-P303        THRU EX-ERRORE-P303
            END-IF
      *
            ADD 1                           TO CTR-CONT-SCARTI
            .
       EX-SCRIVI-SCARTI.
           EXIT.
      *================================================================*
R11817 SCRIVI-SCART2.
      *
            MOVE YPCRPGPF-WORK              TO OUSCART2-PARTE1
            MOVE VG0000R                    TO OUSCART2-PARTE2
            WRITE OUSCART2-REC.
      *
            IF NOT OUS2-NORMAL
               MOVE ST-OUSCART2             TO P303-FILE-STATUS
               MOVE '08'                    TO P303-MSGER-RIF
               MOVE 'OUSCART2'              TO P303-MSGER-FILE
               MOVE 'WRITE'                 TO P303-MSGER-TIPO
               MOVE 'ERRORE WRITE FILE SCA2' TO P303-MSGER-DESCR
               PERFORM ERRORE-P303        THRU EX-ERRORE-P303
            END-IF
      *
            ADD 1                           TO CTR-CONT-SCART2
            .
R11817 EX-SCRIVI-SCART2.
           EXIT.
      *================================================================*
R05818 SCRIVI-FXML.
      *
            WRITE OUTFXML-REC.
      *
            IF NOT OXML-NORMAL
               MOVE ST-OUTFXML              TO P303-FILE-STATUS
               MOVE '18'                    TO P303-MSGER-RIF
               MOVE 'OUTFXML '              TO P303-MSGER-FILE
               MOVE 'WRITE'                 TO P303-MSGER-TIPO
               MOVE 'ERRORE WRITE FILE FXML' TO P303-MSGER-DESCR
               PERFORM ERRORE-P303        THRU EX-ERRORE-P303
            END-IF
      *
            ADD 1                           TO CTR-CONT-FXML
            .
R05818 EX-SCRIVI-FXML.
           EXIT.
      *================================================================*
R11422 SCRIVI-FXM2.
      *
            MOVE  YPCRREQX             TO OUTFXM2-REC.
            WRITE OUTFXM2-REC.
      *
            IF NOT OXM2-NORMAL
               MOVE ST-OUTFXM2              TO P303-FILE-STATUS
               MOVE '18'                    TO P303-MSGER-RIF
               MOVE 'OUTFXM2 '              TO P303-MSGER-FILE
               MOVE 'WRITE'                 TO P303-MSGER-TIPO
               MOVE 'ERRORE WRITE FILE FXM2' TO P303-MSGER-DESCR
               PERFORM ERRORE-P303        THRU EX-ERRORE-P303
            END-IF
      *
            ADD 1                           TO CTR-CONT-FXM2
            .
R11422 EX-SCRIVI-FXM2.
           EXIT.
      *================================================================*
R12019 SCRIVI-BILLCCB.
      *
            WRITE BILLCCB-REC             FROM YPCRBILC-REC
      *
            IF NOT BILC-NORMAL
               MOVE ST-BILLCCB              TO P303-FILE-STATUS
               MOVE '15'                    TO P303-MSGER-RIF
               MOVE 'BILLCCB '              TO P303-MSGER-FILE
               MOVE 'WRITE'                 TO P303-MSGER-TIPO
               MOVE 'ERRORE WRITE FILE BILLCCB'
                                            TO P303-MSGER-DESCR
               PERFORM ERRORE-P303        THRU EX-ERRORE-P303
            END-IF
      *
            ADD 1                           TO CTR-CONT-SCRITTI-BILLCCB
            .
R12019 EX-SCRIVI-BILLCCB.
           EXIT.
      *================================================================*
R14316 SCRIVI-DCD-APER.
      *
      *--* Inizializza area DCD
           INITIALIZE VG0000R
      *
           PERFORM R315-INCREMENTA-PROGR  THRU R315-INCREMENTA-PROGR-EX
           PERFORM IMPO-DCD-DEF           THRU F-IMPO-DCD-DEF
           PERFORM WRITE-OUTDCD           THRU F-WRITE-OUTDCD
      * se il canone e/c � KO devo fare l'apertura anche per il
      * mancato coec
           IF  PGPF-FUNCT-CODE = '200'
           AND PGPF-PAYMT-TYPE = 'PSF'
                MOVE '004900106'         TO VG000-COD-TIP-PART
                MOVE 'A'                 TO VG000-FLG-SGN
                MOVE 'DCBCFACP'          TO VG000-COD-CONT
                PERFORM WRITE-OUTDCD   THRU F-WRITE-OUTDCD
           END-IF.

           .
       EX-SCRIVI-DCD-APER.
           EXIT.
      *================================================================*
R05818 SCRIVI-DCD-APER-CCB.
      *
      *--* Inizializza area DCD
           INITIALIZE VG0000R
      *
           PERFORM R315-INCREMENTA-PROGR  THRU R315-INCREMENTA-PROGR-EX
           INITIALIZE VG0000R
           MOVE 01                          TO VG000-COD-SOC
           MOVE 'SIC49'                     TO VG000-COD-ENTE-4LIV
           MOVE '00'                        TO VG000-COD-UFF
           MOVE 'AP'                        TO VG000-FLG-SEL-TIP-OPE
           MOVE '004900211'                 TO VG000-COD-TIP-PART
           MOVE 'D'                         TO VG000-FLG-SGN
      *
           MOVE 'DCBAMCBS'                  TO VG000-COD-CONT
      *
           MOVE 'SIC49'                     TO VG000-COD-ENTE-4LIV-ORIG
           MOVE '00'                        TO VG000-COD-UFF-ORIG
           MOVE W100-PROGR-APERTURA        TO VG000-CNT-PRG-MOV-PAR
ZANCHI*****MOVE 'P'                        TO VG000-CNT-PRG-MOV-PAR(1:1)
ZANCHI***** nel 1� byte della partita passo il valore alfabetico
ZANCHI***** dell'orario per gestire le 2 fasi OPC in pari data
           PERFORM R320-CONVERTI-ORA      THRU R320-CONVERTI-ORA-EX
           MOVE WS-DATA-AAMMGG              TO VG000-DAT-CONTABILE-AUTO
      *
           MOVE ZEROES                      TO VG000-DAT-SCADENZA
           MOVE ZEROES                      TO VG000-DAT-CHD
           MOVE 'SIC49'                     TO VG000-COD-4LI-NEW-CHS
           MOVE '00'                        TO VG000-COD-UFF-NEW-CHS
           MOVE ZEROES                      TO VG000-DAT-VAL
           MOVE ZEROES                      TO VG000-DAT-CAR
           MOVE ZEROES                      TO VG000-DAT-EMI-EFF
           MOVE ZEROES                      TO VG000-DAT-SCD-EFF
           MOVE 'FULL ACQUIRING - Addebiti merchant su conti bancari'
                                            TO VG000-DES-TT-MOV-PAR
           COMPUTE VG000-IMP-MOV = PGPF-PAYMT-TOT
           MOVE 'EUR'                       TO VG000-COD-DIVISA
           MOVE ZEROES                      TO VG000-IMP-MOV-C
           MOVE WS-DATA-AAAAMMGG            TO VG000-DAT-SOL
           MOVE 'FAC'                       TO VG000-COD-PRV-LAV
           MOVE SPACES                      TO VG000-COD-PRV-SLV
           MOVE WS-DATA-AAAAMMGG            TO VG000-DAT-CONTABILE-X8
      *
           MOVE 'CCBA'                      TO VG000-KEY-PROC(1:4)
           MOVE WS-IBAN                     TO VG000-KEY-PROC(5:27)
R12019*    MOVE WS-AREA-APPO-YPOE-DESC      TO VG000-KEY-PROC(32:)
R12019     MOVE WS-AREA-APPO-YPOE-DESC      TO VG000-KEY-PROC(32:20)
R12019     MOVE YPCWFAMI-O-CODI-FISC        TO VG000-KEY-PROC(52:16)
           PERFORM WRITE-OUTDCD           THRU F-WRITE-OUTDCD
           .
R05818 EX-SCRIVI-DCD-APER-CCB.
           EXIT.
      *==============================================================*
R14316 R315-INCREMENTA-PROGR.
      *
           ADD 1                         TO W100-PROGR
           MOVE 'E'                      TO YPCCP008-TIPO-FUNZ
           MOVE W100-PROGR               TO YPCCP008-D-CAMPO
           CALL W100-ROUT-DECO USING     YPCCP008
           MOVE YPCCP008-E-CAMPO         TO W100-E-CAMPO
           MOVE W100-PROGR-ESA           TO W100-PROGR-APERTURA
           .
       R315-INCREMENTA-PROGR-EX.
           EXIT.
      *==============================================================*
R14316 R320-CONVERTI-ORA.
      *
           EVALUATE WS-ORA-HH
               WHEN '00'
                    MOVE 'A' TO VG000-CNT-PRG-MOV-PAR(1:1)
               WHEN '01'
                    MOVE 'B' TO VG000-CNT-PRG-MOV-PAR(1:1)
               WHEN '02'
                    MOVE 'C' TO VG000-CNT-PRG-MOV-PAR(1:1)
               WHEN '03'
                    MOVE 'D' TO VG000-CNT-PRG-MOV-PAR(1:1)
               WHEN '04'
                    MOVE 'E' TO VG000-CNT-PRG-MOV-PAR(1:1)
               WHEN '05'
                    MOVE 'F' TO VG000-CNT-PRG-MOV-PAR(1:1)
               WHEN '06'
                    MOVE 'G' TO VG000-CNT-PRG-MOV-PAR(1:1)
               WHEN '07'
                    MOVE 'H' TO VG000-CNT-PRG-MOV-PAR(1:1)
               WHEN '08'
                    MOVE 'I' TO VG000-CNT-PRG-MOV-PAR(1:1)
               WHEN '09'
                    MOVE 'J' TO VG000-CNT-PRG-MOV-PAR(1:1)
               WHEN '10'
                    MOVE 'K' TO VG000-CNT-PRG-MOV-PAR(1:1)
               WHEN '11'
                    MOVE 'L' TO VG000-CNT-PRG-MOV-PAR(1:1)
               WHEN '12'
                    MOVE 'M' TO VG000-CNT-PRG-MOV-PAR(1:1)
               WHEN '13'
                    MOVE 'N' TO VG000-CNT-PRG-MOV-PAR(1:1)
               WHEN '14'
                    MOVE 'O' TO VG000-CNT-PRG-MOV-PAR(1:1)
               WHEN '15'
                    MOVE 'P' TO VG000-CNT-PRG-MOV-PAR(1:1)
               WHEN '16'
                    MOVE 'Q' TO VG000-CNT-PRG-MOV-PAR(1:1)
               WHEN '17'
                    MOVE 'R' TO VG000-CNT-PRG-MOV-PAR(1:1)
               WHEN '18'
                    MOVE 'S' TO VG000-CNT-PRG-MOV-PAR(1:1)
               WHEN '19'
                    MOVE 'T' TO VG000-CNT-PRG-MOV-PAR(1:1)
               WHEN '20'
                    MOVE 'U' TO VG000-CNT-PRG-MOV-PAR(1:1)
               WHEN '21'
                    MOVE 'V' TO VG000-CNT-PRG-MOV-PAR(1:1)
               WHEN '22'
                    MOVE 'W' TO VG000-CNT-PRG-MOV-PAR(1:1)
               WHEN '23'
                    MOVE 'X' TO VG000-CNT-PRG-MOV-PAR(1:1)
           END-EVALUATE.
           .
       R320-CONVERTI-ORA-EX.
           EXIT.
      *==============================================================*
R14316 IMPO-DCD-DEF.
      *
           INITIALIZE VG0000R
           MOVE 01                          TO VG000-COD-SOC
           MOVE 'SIC49'                     TO VG000-COD-ENTE-4LIV
           MOVE '00'                        TO VG000-COD-UFF
           MOVE 'AP'                        TO VG000-FLG-SEL-TIP-OPE
      *

           EVALUATE TRUE
               WHEN PGPF-FUNCT-CODE = '300'
                AND PGPF-DB-CR-FLAG = 'C'
                    MOVE '004900210'         TO VG000-COD-TIP-PART
                    MOVE 'A'                 TO VG000-FLG-SGN
                    MOVE 'DCBAKKMP'          TO VG000-COD-CONT
               WHEN PGPF-FUNCT-CODE = '300'
                AND PGPF-DB-CR-FLAG = 'D'
                    MOVE '004900209'         TO VG000-COD-TIP-PART
                    MOVE 'D'                 TO VG000-FLG-SGN
                    MOVE 'DCBAKOMP'          TO VG000-COD-CONT
               WHEN PGPF-FUNCT-CODE     = '200'
                AND PGPF-DB-CR-FLAG     = 'C'
                    MOVE '004900213'         TO VG000-COD-TIP-PART
                    MOVE 'A'                 TO VG000-FLG-SGN
                    MOVE 'DCBAKKCM'          TO VG000-COD-CONT
               WHEN PGPF-FUNCT-CODE     = '200'
                AND PGPF-DB-CR-FLAG     = 'D'
                    MOVE '004900212'         TO VG000-COD-TIP-PART
                    MOVE 'D'                 TO VG000-FLG-SGN
                    MOVE 'DCBAKOCM'          TO VG000-COD-CONT
           END-EVALUATE.
      *
           MOVE 'SIC49'                     TO VG000-COD-ENTE-4LIV-ORIG
           MOVE '00'                        TO VG000-COD-UFF-ORIG
           MOVE W100-PROGR-APERTURA        TO VG000-CNT-PRG-MOV-PAR
ZANCHI*****MOVE 'P'                        TO VG000-CNT-PRG-MOV-PAR(1:1)
ZANCHI***** nel 1� byte della partita passo il valore alfabetico
ZANCHI***** dell'orario per gestire le 2 fasi OPC in pari data
           PERFORM R320-CONVERTI-ORA      THRU R320-CONVERTI-ORA-EX

           MOVE WS-DATA-AAMMGG              TO VG000-DAT-CONTABILE-AUTO
      *
           MOVE ZEROES                      TO VG000-DAT-SCADENZA
           MOVE ZEROES                      TO VG000-DAT-CHD
           MOVE 'SIC49'                     TO VG000-COD-4LI-NEW-CHS
           MOVE '00'                        TO VG000-COD-UFF-NEW-CHS
           MOVE ZEROES                      TO VG000-DAT-VAL
           MOVE ZEROES                      TO VG000-DAT-CAR
           MOVE ZEROES                      TO VG000-DAT-EMI-EFF
           MOVE ZEROES                      TO VG000-DAT-SCD-EFF
           MOVE 'Postepay Evolution Business  KO BATCH FULL ACQUIRING'
                                            TO VG000-DES-TT-MOV-PAR
           COMPUTE VG000-IMP-MOV = PGPF-PAYMT-TOT
           MOVE 'EUR'                       TO VG000-COD-DIVISA
           MOVE ZEROES                      TO VG000-IMP-MOV-C
           MOVE WS-DATA-AAAAMMGG            TO VG000-DAT-SOL
           MOVE 'FAC'                       TO VG000-COD-PRV-LAV
           MOVE SPACES                      TO VG000-COD-PRV-SLV
           MOVE WS-DATA-AAAAMMGG            TO VG000-DAT-CONTABILE-X8
      *
      *
           MOVE 'MM  '                      TO VG000-KEY-PROC(1:4)
           MOVE WS-IBAN                     TO VG000-KEY-PROC(5:27)
           MOVE WS-AREA-APPO-YPOE-DESC      TO VG000-KEY-PROC(32:)
           .
       F-IMPO-DCD-DEF.
           EXIT.
      *================================================================*
      * WRITE FLUSSO DCD                                               *
      *================================================================*
R14316 WRITE-OUTDCD.
      *
           WRITE REC-OUTDCD               FROM VG0000R
      *
           IF NOT OUTDCD-NORMAL
              MOVE ST-OUTDCD                TO P303-FILE-STATUS
              MOVE '09'                     TO P303-MSGER-RIF
              MOVE 'OUTDCD  '               TO P303-MSGER-FILE
              MOVE 'WRITE'                  TO P303-MSGER-TIPO
              MOVE 'ERRORE WRITE FILE OUTDCD ' TO P303-MSGER-DESCR
              PERFORM ERRORE-P303         THRU EX-ERRORE-P303
           END-IF
      *
           ADD     1                        TO CTR-CONT-SCRITTI-DCD
           .
       F-WRITE-OUTDCD.
           EXIT.
      *================================================================*
       SCRIVI-ERRO-MAIL.
R11422*
R11422*--* Se il record } un CC-BANCARI(SDD) andato a buon fine
R11422*--* non scrive il file
R11422     IF WS-ELAB-CC-BANCARI-SI
R11422        GO TO EX-SCRIVI-ERRO-MAIL
R11422     END-IF
      *
R14316*--* Se importo a zero non scrive l'errore
R14316     IF PGPF-PAYMT-TOT = ZEROES
R14316        GO TO EX-SCRIVI-ERRO-MAIL
R14316     END-IF
      *
           IF WS-SCRI-YPOE-NO
              MOVE WS-AREA-TEST-YPOE        TO WS-AREA-APPO-YPOE
              PERFORM WRIT-YPOE           THRU F-WRIT-YPOE
              MOVE  SPACES                  TO WS-AREA-APPO-YPOE
              PERFORM WRIT-YPOE           THRU F-WRIT-YPOE
           END-IF
      *
R14316*--* Scelto di visualizzare IBAN e Importo
R14316     MOVE WS-AREA-APPO-YPOE-DESC      TO WS-AREA-APPO-YPOE-D
R14316*    MOVE PGPF-PAYEMT-UID             TO WS-AREA-APPO-PAYEMT-UID
R11422*    MOVE PGPF-BANK-ACCOUNT           TO WS-AREA-APPO-YPOE-IBAN
R11422     MOVE YPCWFAMI-O-NUME-RAPP-FA     TO WS-AREA-APPO-YPOE-NRFA
R11422     IF PGPF-FUNCT-CODE NOT = '301'
R11422        MOVE PGPF-ME-ID-CODE          TO WS-AREA-APPO-YPOE-MID
R11422     ELSE
R11422        MOVE PGPF-ME-ID-CODE-301      TO WS-AREA-APPO-YPOE-MID
R11422     END-IF
R14316*    MOVE PGPF-BANK-ACC-TYP           TO WS-AREA-APPO-ACC-TYPE
R14316     MOVE PGPF-PAYMT-TOT              TO WS-AREA-APPO-YPOE-IMPO
      *
           PERFORM WRIT-YPOE              THRU F-WRIT-YPOE
           ADD 1                            TO CTR-CONT-SCRITTI-YPOERRO

           MOVE SPACES                      TO WS-AREA-APPO-YPOE
           PERFORM WRIT-YPOE              THRU F-WRIT-YPOE
           .
       EX-SCRIVI-ERRO-MAIL.
           EXIT.
      *================================================================*
       WRIT-YPOE.
      *
           WRITE YPOERRO-REC              FROM WS-AREA-APPO-YPOE

           IF NOT OERR-NORMAL
              PERFORM   IMPO-ERRO-WRIT-YPOE
                 THRU F-IMPO-ERRO-WRIT-YPOE
           END-IF.
      *
           SET WS-SCRI-YPOE-SI              TO TRUE
           .
      *
       F-WRIT-YPOE.
           EXIT.
      *================================================================*
       IMPO-ERRO-WRIT-YPOE.
      *
           MOVE SPACES                      TO YPCWS001-RIGA
           MOVE SPACES                      TO YP-MSGERR
      *
           STRING '-FILE DI OUTPUT: YPOERRO  -ERRORE WRITE -STC: '
                  ST-YPOERRO
                  DELIMITED BY SIZE
                  INTO YP-MSGERR-1
           END-STRING
      *
           PERFORM GEST-ERRO-SU-TRE-RIGH  THRU F-GEST-ERRO-SU-TRE-RIGH
           .
       F-IMPO-ERRO-WRIT-YPOE.
           EXIT.
      *================================================================*
       IMPOSTA-DATI-OUT.
      *
           INITIALIZE CRVSD50-RECF
      *
           MOVE '00001'                    TO CRVSD50-ISTITUT
           MOVE 1                          TO CTR-PROGRES
           MOVE CTR-PROGRES                TO CRVSD50-PROGRES
           MOVE '1'                        TO CRVSD50-TRAGGRU
           MOVE 'CC '                      TO CRVSD50-TIPSERV
           MOVE COMSD50-FILIALE            TO CRVSD50-FILIALE
           MOVE COMSD50-RAPPORT            TO CRVSD50-RAPPORT
           MOVE COMSD50-CATRAPP            TO CRVSD50-CATRAPP
010715*--* La data contabile deve essere uguale a quella solare
220715*--* La data contabile NON deve essere uguale a quella solare
220715     MOVE PGPF-PAYMT-GEN-DATE        TO CRVSD50-DATCONT
010715*    MOVE WS-DATA-AAAAMMGG           TO CRVSD50-DATCONT
           MOVE '01'                       TO CRVSD50-TIPOINF
           MOVE W100-CAUSALE               TO CRVSD50-CAUSALE
           IF PGPF-DB-CR-FLAG = 'C'
              MOVE  'A'                    TO CRVSD50-FLAGDA
           ELSE
              MOVE  'D'                    TO CRVSD50-FLAGDA
           END-IF
           COMPUTE CRVSD50-IMPOPER =  PGPF-PAYMT-TOT / 100
R14316     IF  (PGPF-FUNCT-CODE = 200)
R14316     AND (PGPF-REFER-PERIOD-DATE-200 NOT = SPACES AND LOW-VALUE)
180112     AND  PGPF-REFER-PERIOD-DATE-200 NOT = 00010101
R14316        MOVE PGPF-REFER-PERIOD-DATE-200 TO CRVSD50-VALLIQU
R14316     ELSE
ZANCHI*    Se commissioni mensili uso la valuta PGPF Altrimenti
ZANCHI*    data pagamento pos + 1 giorno
ZANCHI        IF PGPF-PAYMT-CYCL = 'M'
ZANCHI           MOVE PGPF-VALUE-DATE TO WK-DATA-NUM
ZANCHI        ELSE
ZANCHI           PERFORM CALC-DATA-PIU-UNO
ZANCHI            THRU F-CALC-DATA-PIU-UNO
ZANCHI        END-IF
R04417        MOVE WK-DATA-NUM                TO CRVSD50-VALLIQU
R14316     END-IF
      *
           MOVE 'EUR'                      TO CRVSD50-DIVISA
           MOVE W100-CODOPE                TO CRVSD50-CODOPE
           COMPUTE CRVSD50-CVALORE = PGPF-PAYMT-TOT / 100
           IF PGPF-BRAND-CODE > SPACE
              MOVE PGPF-BRAND-CODE            TO W100-BRAND-CODE
                                                 W100-BRAND-CODE-CH
              MOVE ' Brand '                  TO W100-DESC-BRAND
                                                 W100-DESC-BRAND-CH
           ELSE
              MOVE SPACE                      TO W100-BRAND-CODE
                                                 W100-BRAND-CODE-CH
              MOVE SPACE                      TO W100-DESC-BRAND
                                                 W100-DESC-BRAND-CH
           END-IF
170511     IF PGPF-NICKNAME > SPACE
170511        MOVE PGPF-NICKNAME              TO W100-NICKNAME
170511                                           W100-NICKNAME-CH
170511     ELSE
170511        MOVE SPACE                      TO W100-NICKNAME
170511                                           W100-NICKNAME-CH
170511     END-IF
           STRING PGPF-BUSINESS-DATE(7:2) '-'
                  PGPF-BUSINESS-DATE(5:2) '-'
                  PGPF-BUSINESS-DATE(1:4)
                  DELIMITED BY SIZE
                  INTO W100-DATA-OPERAZ
           END-STRING
160408     MOVE PGPF-ME-ID-CODE           TO W100-MERCHANT
160408                                       W100-MERCHANT-CH
           IF PGPF-FUNCT-CODE   = 200
              MOVE W100-DESCMOV-200        TO CRVSD50-DESCMOV
            ELSE
              MOVE W100-DESCMOV            TO CRVSD50-DESCMOV
           END-IF
R03118*--* In data 19/03/18, CC non � pronta con la copy D50 aggiornata
R03118*    MOVE PGPF-PAYEMT-UID            TO CRVSD50-INFOPFM-SS
R03118     MOVE PGPF-PAYEMT-UID            TO CRVSD50-AREAPAS(159:18)
R10421     IF PGPF-DB-CR-FLAG = 'C'
R10421        MOVE  'A'                    TO CRVSD50-AREAPAS(177:01)
R10421     ELSE
R10421        MOVE  'D'                    TO CRVSD50-AREAPAS(177:01)
R10421     END-IF
      *
            .
       EX-IMPOSTA-DATI-OUT.
           EXIT.
      *================================================================*
R05818 IMPOSTA-DATI-OUT-CC-VTPIE.
      *
      *
           PERFORM   SELE-TABE-GEP-CCB
              THRU F-SELE-TABE-GEP-CCB
      *
           PERFORM RICERCA-FRAZIONARIO-CC
              THRU EX-RICERCA-FRAZIONARIO-CC
      *
           INITIALIZE CRVSD50-RECF
      *
           MOVE '00001'                    TO CRVSD50-ISTITUT
           MOVE 1                          TO CTR-PROGRES
           MOVE CTR-PROGRES                TO CRVSD50-PROGRES
           MOVE '1'                        TO CRVSD50-TRAGGRU
           MOVE 'CC '                      TO CRVSD50-TIPSERV
           MOVE COMSD50-FILIALE            TO CRVSD50-FILIALE
           MOVE COMSD50-RAPPORT            TO CRVSD50-RAPPORT
           MOVE COMSD50-CATRAPP            TO CRVSD50-CATRAPP
*******  imposto data valuta e data contabile con la data di sistema
           MOVE WS-DATA-AAAAMMGG           TO WS-PAYMT-GEN-DATE
           MOVE WS-PAYMT-GEN-DATE          TO CRVSD50-DATCONT
           MOVE CRVSD50-DATCONT            TO CRVSD50-VALLIQU
           MOVE '01'                       TO CRVSD50-TIPOINF
           MOVE  YPCRTCCB-CAUSALE          TO CRVSD50-CAUSALE
           MOVE  'A'                       TO CRVSD50-FLAGDA
           MOVE   WS-TOTALI-IMPO-CC        TO CRVSD50-IMPOPER
R05818*    MOVE  WS-VALUE-DATE  TO PGPF-VALUE-DATE
R05818*    PERFORM CALC-DATA-PIU-UNO
R05818*      THRU F-CALC-DATA-PIU-UNO
R05818*    MOVE WK-DATA-NUM                TO CRVSD50-VALLIQU
      *
           MOVE 'EUR'                      TO CRVSD50-DIVISA
           MOVE YPCRTCCB-CODOPE            TO CRVSD50-CODOPE
           MOVE WS-TOTALI-IMPO-CC          TO CRVSD50-CVALORE
           STRING WS-BUSINESS-DATE(7:2) '-'
                  WS-BUSINESS-DATE(5:2) '-'
                  WS-BUSINESS-DATE(1:4)
                  DELIMITED BY SIZE
                  INTO W100-DATA-OPER
           END-STRING
           MOVE W100-DESCMOV-CC-VTPIE   TO CRVSD50-DESCMOV
*********  MOVE PGPF-PAYEMT-UID         TO CRVSD50-AREAPAS(159:18)
      *
            .
R05818 EX-IMPOSTA-DATI-OUT-CC-VTPIE.
           EXIT.
      *================================================================*
R12019 IMPOSTA-BILLCCB.
      *
           INITIALIZE YPCRBILC-REC
           MOVE ZERO               TO WK-IMPO-BILC
           MOVE ZERO               TO WK-IMPO-BILC-NUME
      *
           MOVE  'PPAY'            TO  YPCRBILC-SOCIETA
           MOVE  YPCWFAMI-O-CONT-BILLI-CCB
                                   TO  YPCRBILC-CONTRATTO-PROVIDER
      *--* Imposta importo con la virgola
           COMPUTE WK-IMPO-BILC-NUME   = PGPF-PAYMT-TOT / 100
           MOVE WK-IMPO-BILC-NUME  TO WK-IMPO-BILC
           MOVE WK-IMPO-BILC       TO YPCRBILC-IMPORTO
      *
           IF PGPF-DB-CR-FLAG = 'C'
              MOVE  'H'            TO  YPCRBILC-DARE-AVERE
           ELSE
              MOVE  'D'            TO  YPCRBILC-DARE-AVERE
           END-IF
           MOVE  'L1'              TO  YPCRBILC-CODICE-IVA
           MOVE  'ASF'             TO  YPCRBILC-PRODOTTO
           MOVE  'N'               TO  YPCRBILC-OMOLOGAZIONE
           MOVE  'N'               TO  YPCRBILC-UNIVERSALE-NONUNIVERS
           .
       EX-IMPOSTA-BILLCCB.
R12019     EXIT.
      *================================================================*
TK1274 IMPOSTA-BILLCCB-2.
      *
DBG==>*    display 'imposta-BILLCCB-2 '
DBG==>*    display 'WK-TROVATO-SU-FPR 1 : ' WK-TROVATO-SU-FPR
DBG==>*    display 'WK-TROVATO-SU-FPR : ' WK-TROVATO-SU-FPR
           PERFORM IMPOSTA-BILLCCB-ASI
              THRU EX-IMPOSTA-BILLCCB-ASI
           PERFORM SCRIVI-BILLCCB
              THRU EX-SCRIVI-BILLCCB
           PERFORM IMPOSTA-BILLCCB-ASA
             THRU EX-IMPOSTA-BILLCCB-ASA
           PERFORM SCRIVI-BILLCCB
              THRU EX-SCRIVI-BILLCCB
              .
      *
           .
TK1274 EX-IMPOSTA-BILLCCB-2.
           EXIT.
      *================================================================*
TK1274 IMPOSTA-BILLCCB-ASI.
      *
           INITIALIZE YPCRBILC-REC
           MOVE ZERO               TO WK-IMPO-BILC
           MOVE ZERO               TO WK-IMPO-BILC-NUME
      *
           MOVE  'PPAY'            TO  YPCRBILC-SOCIETA
           MOVE  YPCWFAMI-O-CONT-BILLI-CCB
                                   TO  YPCRBILC-CONTRATTO-PROVIDER
      *--* Imposta importo con la virgola
           COMPUTE WK-IMPO-BILC-NUME ROUNDED =
           (PGPF-PAYMT-TOT * TFPR-PERCENTUALE-POSTE) / 1000000

           MOVE WK-IMPO-BILC-NUME  TO WK-IMPO-BILC-NUME-ASI
           MOVE WK-IMPO-BILC-NUME  TO WK-IMPO-BILC
           MOVE WK-IMPO-BILC       TO YPCRBILC-IMPORTO
      *
           IF PGPF-DB-CR-FLAG = 'C'
              MOVE  'H'            TO  YPCRBILC-DARE-AVERE
           ELSE
              MOVE  'D'            TO  YPCRBILC-DARE-AVERE
           END-IF
           MOVE  'L1'              TO  YPCRBILC-CODICE-IVA
           MOVE  'ASI'             TO  YPCRBILC-PRODOTTO
           MOVE  'N'               TO  YPCRBILC-OMOLOGAZIONE
           MOVE  'N'               TO  YPCRBILC-UNIVERSALE-NONUNIVERS
           .
TK1274 EX-IMPOSTA-BILLCCB-ASI.
TK1274     EXIT.
      *================================================================*
TK1274 IMPOSTA-BILLCCB-ASA.
      *
           INITIALIZE YPCRBILC-REC
           MOVE ZERO               TO WK-IMPO-BILC
           MOVE ZERO               TO WK-IMPO-BILC-NUME
      *
           MOVE  'PPAY'            TO  YPCRBILC-SOCIETA
           MOVE  YPCWFAMI-O-CONT-BILLI-CCB
                                   TO  YPCRBILC-CONTRATTO-PROVIDER
      *--* Imposta importo con la virgola
           COMPUTE WK-IMPO-BILC-NUME  ROUNDED =
                (PGPF-PAYMT-TOT / 100 ) - WK-IMPO-BILC-NUME-ASI

           MOVE WK-IMPO-BILC-NUME  TO WK-IMPO-BILC
           MOVE WK-IMPO-BILC       TO YPCRBILC-IMPORTO
      *
           IF PGPF-DB-CR-FLAG = 'C'
              MOVE  'H'            TO  YPCRBILC-DARE-AVERE
           ELSE
              MOVE  'D'            TO  YPCRBILC-DARE-AVERE
           END-IF
           MOVE  'L1'              TO  YPCRBILC-CODICE-IVA
           MOVE  'ASA'             TO  YPCRBILC-PRODOTTO
           MOVE  'N'               TO  YPCRBILC-OMOLOGAZIONE
           MOVE  'N'               TO  YPCRBILC-UNIVERSALE-NONUNIVERS
           .
TK1274 EX-IMPOSTA-BILLCCB-ASA.
TK1274     EXIT.
      *================================================================*
R05818 IMPOSTA-DATI-FXML.
      *
           INITIALIZE OUTFXML-DATI
           INITIALIZE WS-DATI-XML
           MOVE WS-IBAN                 TO  WS-IBAN-DEST
           COMPUTE WS-IMPO-MOV = PGPF-PAYMT-TOT / 100
           MOVE PGPF-PAYEMT-UID        TO   WS-PAYMENT-UID
      *
R13519*--* Nel caso il rapporto non esista, in base dati del FA,
R13519*--* nel campo ragione sociale viene impostata il valore
R13519*--* del campo insegna, passato da SIA nel flusso PGPF
R13519     IF YPCWFAMI-O-NUME-RAPP-FA = SPACES OR LOW-VALUE
R13519        IF PGPF-ACCT-OWNER-NAM  > SPACE
R13519           MOVE PGPF-ACCT-OWNER-NAM  TO WS-RAGI-SOC
R13519        ELSE
R13519           MOVE PGPF-NAME            TO WS-RAGI-SOC
R13519        END-IF
R13519     ELSE
              PERFORM CHIAMA-ANAGRAFE
              THRU EX-CHIAMA-ANAGRAFE
R13519     END-IF
      *
           IF PGPF-BRAND-CODE > SPACE
              MOVE PGPF-BRAND-CODE            TO W100-BRAND-CODE
                                                 W100-BRAND-CODE-CH
              MOVE ' Brand '                  TO W100-DESC-BRAND
                                                 W100-DESC-BRAND-CH
           ELSE
              MOVE SPACE                      TO W100-BRAND-CODE
                                                 W100-BRAND-CODE-CH
              MOVE SPACE                      TO W100-DESC-BRAND
                                                 W100-DESC-BRAND-CH
           END-IF
           IF PGPF-NICKNAME > SPACE
              MOVE PGPF-NICKNAME              TO W100-NICKNAME
                                                 W100-NICKNAME-CH
           ELSE
              MOVE SPACE                      TO W100-NICKNAME
                                                 W100-NICKNAME-CH
           END-IF
           STRING PGPF-BUSINESS-DATE(7:2) '-'
                  PGPF-BUSINESS-DATE(5:2) '-'
                  PGPF-BUSINESS-DATE(1:4)
                  DELIMITED BY SIZE
                  INTO W100-DATA-OPERAZ
           END-STRING
           MOVE PGPF-ME-ID-CODE           TO W100-MERCHANT
                                             W100-MERCHANT-CH
           IF PGPF-FUNCT-CODE   = 200
              MOVE W100-DESCMOV-200        TO WS-DESC-MOV
            ELSE
              MOVE W100-DESCMOV            TO WS-DESC-MOV
           END-IF
           MOVE  WS-DESC-MOV               TO OUTFXML-DESC
      *
           MOVE  '1'          TO OUTFXML-TIPO-REC
           MOVE WS-DATI-XML   TO OUTFXML-DATI.
      *
R05818 EX-IMPOSTA-DATI-FXML.
           EXIT.
      *================================================================*
R05818 IMPOSTA-DATI-T-FXML.
      *
           INITIALIZE OUTFXML-DATI
      *
           MOVE  YPCRTCCB-IBAN     TO  WS-IBAN-MITT
           MOVE WS-NUM-OPER-CC TO WS-NUM-TRAN
           MOVE WS-TOTALI-IMPO-CC TO WS-TOTALE-IMP
           MOVE  '0'           TO OUTFXML-TIPO-REC
           MOVE WS-DATI-XML-CODA      TO OUTFXML-DATI.
           MOVE  SPACES               TO OUTFXML-DESC.
      *
R05818 EX-IMPOSTA-DATI-T-FXML.
           EXIT.
      *================================================================*
R05818 ADD-TOTALI.
      *
           COMPUTE WS-TOTALI-IMPO-CC =
                           WS-TOTALI-IMPO-CC + (PGPF-PAYMT-TOT / 100)
           COMPUTE WS-NUM-OPER-CC = WS-NUM-OPER-CC + 1
      *
            .
R05818 EX-ADD-TOTALI.
           EXIT.
      *================================================================*
R05818 CHIAMA-ANAGRAFE.
      *
           INITIALIZE    AREA-ACS108A
           INITIALIZE                        L-ACS108-ARG
           MOVE ZEROES                    TO L-ACS108-I-BANCA
           MOVE ' '                       TO L-ACS108-I-TIPO-RICH
           MOVE ZEROES                    TO L-ACS108-I-DATA-RIF
*********  MOVE 'CC '                     TO L-ACS108-I-SERVIZIO
           MOVE 'FA '                     TO L-ACS108-I-SERVIZIO
           MOVE ZEROES                    TO WK-END
*********  MOVE WS-IBAN-DEST(16:12)          TO L-ACS108-I-NUMERO-X
           MOVE YPCWFAMI-O-NUME-RAPP-FA  TO L-ACS108-I-NUMERO
           MOVE YPCWFAMI-O-FILI-RAPP     TO L-ACS108-I-FILIALE
           CALL WK-ACS108BT             USING L-ACS108-ARG.
           .
           IF  L-ACS108-RET-CODE = ZEROES
               IF    L-ACS108-COGNOME  = SPACES AND
                     L-ACS108-NOME  = SPACES
                    MOVE  L-ACS108-RAGSOC-1      TO WS-RAGI-SOC
                    MOVE  L-ACS108-RAGSOC-1      TO WS-RAGI-SOC
               ELSE
                   MOVE  L-ACS108-COGNOME        TO WS-RAGI-SOC(1:30)
                   MOVE  L-ACS108-NOME           TO WS-RAGI-SOC(31:)
               END-IF
               MOVE  L-ACS108-IND-SEDE-LEG   TO WS-INDIRIZZO
               MOVE  L-ACS108-CAP-SEDE-LEG   TO WS-CAP
               MOVE  L-ACS108-LOC-SEDE-LEG   TO WS-LOC
               MOVE  L-ACS108-PROV-SEDE-LEG  TO WS-PROV
               MOVE  L-ACS108-NAZ-SEDE-LEG   TO WS-NAZ
R13621     ELSE
R13621         DISPLAY '***** ERRORE DA ANAGRAFE ****** '
R13621         DISPLAY 'L-ACS108-RET-CODE : '  L-ACS108-RET-CODE
R13621         DISPLAY 'L-ACS108-I-NUMERO : ' L-ACS108-I-NUMERO
R13621         DISPLAY 'L-ACS108-I-FILIALE: ' L-ACS108-I-FILIALE
R13621         IF PGPF-ACCT-OWNER-NAM  > SPACE
R13621            MOVE PGPF-ACCT-OWNER-NAM  TO WS-RAGI-SOC
R13621         ELSE
R13621            MOVE PGPF-NAME            TO WS-RAGI-SOC
R13621         END-IF
R13621     END-IF
            .
      *
R05818 EX-CHIAMA-ANAGRAFE.
           EXIT.
      *================================================================*
R04417 CALC-DATA-PIU-UNO.
      *
      * CRVSD50-VALLIQU = DATA PAGAMENTO + 1
      *
           MOVE SPACES                    TO UTDATA-PARAM.
           INITIALIZE UTDATA-PARAM.
      *
ZANCHI*    Se commissioni mensili uso la vauta PGPF Altrimenti
ZANCHI*    business date + 1
ZANCHI     MOVE PGPF-BUSINESS-DATE        TO UTDATA-DATA-1.
           MOVE 3                         TO UTDATA-FUNZIONE.
           MOVE 1                         TO UTDATA-GIORNI.

           CALL WS-XSCDAT USING UTDATA-PARAM

           IF UTDATA-ERRORE NOT = ZEROES
              DISPLAY '*** ERRORE CALCOLO ROUTINE DATA XSCDAT *** '
              DISPLAY '  COD-ERR ---> '  UTDATA-ERRORE
              DISPLAY '  DATA    ---> '  UTDATA-DATA-1
              DISPLAY '  IMPOSTATA DATA VALUTA = DATA ELABORAZIONE'
              MOVE PGPF-VALUE-DATE        TO WK-DATA-NUM
           ELSE
              MOVE UTDATA-DATA-2          TO WK-DATA-NUM
           END-IF
           .
       F-CALC-DATA-PIU-UNO.
           EXIT.
      *================================================================*
R14316 IMPOSTA-DATI-OUT-CONT.
      *
           INITIALIZE                       XYCRCONT
      *
           MOVE 'US'                     TO XYCRCONT-COMPAGNIA
           MOVE W100-CODOPE              TO XYCRCONT-TIPO-MSG
           IF PGPF-DB-CR-FLAG = 'C'
              MOVE  'A'                  TO XYCRCONT-SEGNO
           ELSE
              MOVE  'D'                  TO XYCRCONT-SEGNO
           END-IF
           MOVE '00'                     TO XYCRCONT-FLAG-OLI
           MOVE SPACES                   TO XYCRCONT-NUMMOVI-POSTE
           MOVE '07601'                  TO XYCRCONT-COD-ABI
           MOVE ZERO                     TO XYCRCONT-COD-GRUPPO
           MOVE WS-PAN-EUC               TO XYCRCONT-PAN-2TRA
           MOVE SPACES                   TO XYCRCONT-PAN-SUPERSIM
           MOVE SPACES                   TO XYCRCONT-041-TERM-ID
           MOVE SPACES                   TO XYCRCONT-RRN
      *
      *    MOVE Z3CLGE90-CLE-PAN-III     TO XYCRCONT-PAN-BCM
           INSPECT Z3CLGE90-CLE-PAN-III TALLYING WK-END
               FOR CHARACTERS BEFORE INITIAL SPACE
           .
           IF WK-END > 17
              MOVE 17 TO WK-END
           .
           MOVE Z3CLGE90-CLE-PAN-III(1:WK-END) TO XYCRCONT-PAN-BCM
      *
           MOVE ZEROES                    TO WK-END
           MOVE Z3CLGE90-CLE-TIPO-CARTA   TO XYCRCONT-TIPO-CARTA
           INSPECT Z3CLGE90-CLE-CONTO TALLYING WK-END
               FOR CHARACTERS BEFORE INITIAL SPACE
           .
           MOVE Z3CLGE90-CLE-CONTO(1:WK-END) TO XYCRCONT-CC-CARTA
           MOVE ZEROES                    TO WK-END
           STRING  '0' Z3CLGE90-CLE-FILIALE DELIMITED BY SIZE INTO
                                             XYCRCONT-FILIALE
           MOVE Z3CLGE90-CLE-CATEG-CONTO  TO XYCRCONT-CATEGORIA
           INSPECT Z3CLGE90-CLE-NDG TALLYING WK-END
               FOR CHARACTERS BEFORE INITIAL SPACE
           .
           MOVE Z3CLGE90-CLE-NDG(1:WK-END) TO XYCRCONT-NDG
           MOVE ZEROES                    TO WK-END
      *
      *    MOVE Z3CLGE90-CLE-TIPO-CARTA  TO XYCRCONT-TIPO-CARTA
      *    MOVE ZERO                     TO XYCRCONT-CC-CARTA
      *    MOVE Z3CLGE90-CLE-FILIALE     TO XYCRCONT-FILIALE
      *    MOVE Z3CLGE90-CLE-CATEG-CONTO TO XYCRCONT-CATEGORIA
      *    MOVE Z3CLGE90-CLE-NDG         TO XYCRCONT-NDG
           MOVE PGPF-PAYMT-TOT           TO XYCRCONT-IMPORTO
           MOVE 2                        TO XYCRCONT-NR-DECIMALI-IMP
      *
R14316     IF  (PGPF-FUNCT-CODE = 200)
R14316     AND (PGPF-REFER-PERIOD-DATE-200 NOT = SPACES AND LOW-VALUE)
R14316       MOVE WS-DATA-AAAAMMGG           TO XYCRCONT-DATA-OPERAZIONE
R14316     ELSE
             MOVE PGPF-BUSINESS-DATE(1:8)    TO XYCRCONT-DATA-OPERAZIONE
R14316     END-IF
      *
           MOVE ZERO                     TO XYCRCONT-TIME-OPERAZIONE
           MOVE W100-CAUSALE             TO XYCRCONT-CAUSALE-INTERNA
           MOVE ZERO                     TO XYCRCONT-BUSINESS-CODE
           MOVE SPACES                   TO XYCRCONT-APPROVAL-CODE
      *-
      *--* Imposta dati per descrizione
           IF PGPF-BRAND-CODE > SPACE
              MOVE PGPF-BRAND-CODE       TO W100-DESPPAY-BR
                                            W100-BRAND-CODE-PPAY
              MOVE ' Brand '             TO W100-DESPPAY-DESC-BR
                                            W100-DESC-BRAND-PPAY
           ELSE
              MOVE SPACE                 TO W100-DESPPAY-BR
                                            W100-BRAND-CODE-PPAY
                                            W100-DESPPAY-DESC-BR
                                            W100-DESC-BRAND-PPAY
           END-IF
           STRING XYCRCONT-DATA-OPERAZIONE(7:2) '-'
                  XYCRCONT-DATA-OPERAZIONE(5:2) '-'
                  XYCRCONT-DATA-OPERAZIONE(1:4)
                  DELIMITED BY SIZE
                  INTO W100-DESPPAY-DT
           END-STRING
           MOVE PGPF-ME-ID-CODE          TO W100-DESPPAY-ME
                                            W100-MERCHANT-PPAY
           IF PGPF-FUNCT-CODE   = 200
              MOVE W100-DESPPAY-200      TO XYCRCONT-ESERCENTE
           ELSE
              MOVE W100-DESPPAY          TO XYCRCONT-ESERCENTE
           END-IF
      *
           MOVE ZERO                     TO XYCRCONT-DATA-OLI
           MOVE ZERO                     TO XYCRCONT-TIME-OLI
           MOVE ZERO                     TO XYCRCONT-IMPORTO-OLI
           MOVE SPACES                   TO XYCRCONT-DIV-ORIGINARIA
           MOVE ZERO                     TO XYCRCONT-IMP-ORIGINARIA
           MOVE ZERO                     TO XYCRCONT-COMMISSIONI
           MOVE 3                        TO XYCRCONT-NR-DECIMALI-COMM
           MOVE ZERO                     TO XYCRCONT-STAN
           .
           MOVE 'R'                      TO XYCRCONT-FLAG-RISCHIO
           MOVE 'MM'                     TO XYCRCONT-TIPO-CONTO
           MOVE ZERO                     TO XYCRCONT-PROG-FLUSSO
????****   MOVE YPDCPGPF-DATE-VALUE      TO XYCRCONT-DATA-FLUSSO
????       MOVE COM-DATE-TIME-N(1:8)     TO XYCRCONT-DATA-FLUSSO
           .
           MOVE Z3CLGE90-CLE-TIPO-DISP   TO XYCRCONT-TIPO-DISPOSITIVA
????  *    MOVE ?????????                TO XYCRCONT-GRUPPO-ESERCENTI
           MOVE Z3CLGE90-CLE-COD-PROD    TO XYCRCONT-COD-PROD
????       MOVE SPACES                   TO XYCRCONT-ARN
????       MOVE SPACES                   TO XYCRCONT-022-POS-DCD
           MOVE PGPF-ME-ID-CODE          TO XYCRCONT-COD-CONV
           MOVE ZERO                     TO XYCRCONT-IMP-FEES
           MOVE SPACES                   TO XYCRCONT-TIPO-COD-ASS
           MOVE '09509'                  TO XYCRCONT-COD-ACQUIRER
           MOVE '12928'                  TO XYCRCONT-033-FORW-INST-ID
R05818     MOVE PGPF-PAYEMT-UID          TO XYCRCONT-PAYEMT-UID
R14217*
R14217     PERFORM TROVA-KEY-RANDOM    THRU F-TROVA-KEY-RANDOM
R14217     MOVE HV-KEY-RANDOM            TO XYCRCONT-KEY-RANDOM-NUM
           .
      *
       EX-IMPOSTA-DATI-OUT-CONT.
           EXIT.
      *================================================================*
       SCRIVI-REC-OUT.
      *
            WRITE  OPECONT-REC            FROM  CRVSD50-RECF
      *
            IF NOT OPEC-NORMAL
               MOVE ST-OPECONT              TO P303-FILE-STATUS
               MOVE '09'                    TO P303-MSGER-RIF
               MOVE 'OPECONT '              TO P303-MSGER-FILE
               MOVE 'WRITE'                 TO P303-MSGER-TIPO
               MOVE 'ERRORE WRITE FILE OPECONT' TO P303-MSGER-DESCR
               PERFORM ERRORE-P303        THRU EX-ERRORE-P303
            END-IF
      *
            ADD     1                       TO CTR-CONT-SCRITTI
            .
       EX-SCRIVI-REC-OUT.
           EXIT.
      *================================================================*
FIANNH SCRIVI-REC-OUTB.
FIANNH*
FIANNH      WRITE  OPECONTB-REC            FROM  CRVSD50-RECF
FIANNH*
FIANNH      IF NOT OPECB-NORMAL
FIANNH         MOVE ST-OPECONTB              TO P303-FILE-STATUS
FIANNH         MOVE '09'                    TO P303-MSGER-RIF
FIANNH         MOVE 'OPECONTB '              TO P303-MSGER-FILE
FIANNH         MOVE 'WRITE'                 TO P303-MSGER-TIPO
FIANNH         MOVE 'ERRORE WRITE FILE OPECONTB' TO P303-MSGER-DESCR
FIANNH         PERFORM ERRORE-P303        THRU EX-ERRORE-P303
FIANNH      END-IF
FIANNH*
FIANNH      ADD     1                       TO CTR-CONT-SCRITTI-B
FIANNH      .
FIANNH EX-SCRIVI-REC-OUTB.
FIANNH     EXIT.
      *================================================================*
R14316 SCRIVI-REC-OUT-CONT.
      *
           WRITE XYDCONT-REC            FROM XYCRCONT
           .
            IF NOT XYDCONT-NORMAL
               MOVE ST-XYDCONT              TO P303-FILE-STATUS
               MOVE '17'                    TO P303-MSGER-RIF
               MOVE 'XYDCONT '              TO P303-MSGER-FILE
               MOVE 'WRITE'                 TO P303-MSGER-TIPO
               MOVE 'ERRORE WRITE FILE XYDCONT' TO P303-MSGER-DESCR
               PERFORM ERRORE-P303        THRU EX-ERRORE-P303
            END-IF
      *
            ADD     1                       TO CTR-CONT-SCRITTI-CONT
            .
       EX-SCRIVI-REC-OUT-CONT.
           EXIT.
      *================================================================*
       INSE-TABE-PGPF.
      *
DBG==>*    DISPLAY 'PGPF-FUNCT-CODE (' PGPF-FUNCT-CODE')'
R08421     IF PGPF-FUNCT-CODE NOT = '301'
DBG==>*       DISPLAY 'ESEGUO IMPO-PGPF'
              PERFORM IMPO-PGPF
              THRU  F-IMPO-PGPF
R08421     ELSE
DBG==>*       DISPLAY 'ESEGUO IMPO-PGPF-301'
R08421        PERFORM IMPO-PGPF-301
R08421        THRU  F-IMPO-PGPF-301
R08421     END-IF
      *
           PERFORM INSE-PGPF
           THRU  F-INSE-PGPF
           .
       F-INSE-TABE-PGPF.
           EXIT.
      *================================================================*
       IMPO-PGPF.
      *
           INITIALIZE DCLYPTBPGPF.
           STRING COM-DATE-TIME-H(1:4) '-'
                  COM-DATE-TIME-H(5:2) '-'
                  COM-DATE-TIME-H(7:2)
                  DELIMITED BY SIZE
                  INTO  YPDCPGPF-PGPF-DATE
           END-STRING.
           STRING COM-DATE-TIME-H(9:2)  ':'
                  COM-DATE-TIME-H(11:2) ':'
                  COM-DATE-TIME-H(13:2)
                  DELIMITED BY SIZE
                  INTO  YPDCPGPF-PGPF-TIME
           END-STRING.
R08421     MOVE SPACE                TO YPDCPGPF-SUMMARY-UID
           MOVE PGPF-FUNCT-CODE      TO YPDCPGPF-FUNCTION-CODE
           MOVE PGPF-PAYEMT-UID      TO YPDCPGPF-PAYMENT-UID
           MOVE PGPF-SRC-COD-IND     TO YPDCPGPF-SOURCE-CODE-IND
           MOVE PGPF-BRAND-CODE      TO YPDCPGPF-BRAND
           MOVE PGPF-BANK-ACC-TYP    TO YPDCPGPF-BANK-ACCOUNT-TYPE
           MOVE PGPF-PAYMT-TYPE      TO YPDCPGPF-PAYMENT-TYPE
           MOVE PGPF-PAYMT-CYCL      TO YPDCPGPF-PAYMENT-CYCLE
           MOVE PGPF-DB-CR-FLAG      TO YPDCPGPF-D-C-IND
           COMPUTE YPDCPGPF-PAYMENT-AMOUNT =
                   PGPF-PAYMT-TOT / 100
           MOVE PGPF-BILL-CURR-CODE  TO YPDCPGPF-BILLING-CURRENCY
           MOVE PGPF-BANK-ACCOUNT    TO YPDCPGPF-BANK-ACCOUNT-NUM
           MOVE PGPF-LEVEL-PAY-CODE  TO YPDCPGPF-LEVEL-PAYMENT-CODE
           MOVE PGPF-ME-ID-CODE      TO YPDCPGPF-MERCHANT-ID
180112*****MOVE PGPF-NICKNAME        TO YPDCPGPF-NOME-MERCHANT
180112     MOVE PGPF-NAME            TO YPDCPGPF-NOME-MERCHANT
           MOVE PGPF-BILLING-FLAG    TO YPDCPGPF-NET-GROSS-IND
           MOVE PGPF-SUPPRESS-CODE   TO YPDCPGPF-SUPPRESS-CODE
R05316     IF PGPF-FUNCT-CODE = '300'
             MOVE PGPF-PAYMENT-SCHEME  TO YPDCPGPF-PAYMENT-SCHEME
             MOVE PGPF-VAT-NUMBER      TO YPDCPGPF-VAT-NUM
             MOVE '0001-01-01'         TO YPDCPGPF-REF-PERIOD-DATE
R05316     END-IF
R05316     IF PGPF-FUNCT-CODE = '200'
R05316       MOVE PGPF-PAYMENT-SCHEME-200  TO YPDCPGPF-PAYMENT-SCHEME
R05316       MOVE PGPF-VAT-NUMBER-200      TO YPDCPGPF-VAT-NUM
R54824       MOVE PGPF-PAYM-REAS-200       TO YPDCPGPF-PAYMENT-REASON
R05316       COMPUTE YPDCPGPF-UNIT-AMOUNT =
R05316               PGPF-UNIT-AMOUNT-DEVICE-200 / 100
R05316       MOVE PGPF-DEVIS-TOTAL-NUMBER-200 TO YPDCPGPF-DEVICE-TOT
R05316       MOVE PGPF-REFER-PERIOD-DATE-200  TO COM-DATE-N
R05316       STRING COM-DATE(1:4) '-'
R05316              COM-DATE(5:2) '-'
R05316              COM-DATE(7:2)
R05316            DELIMITED BY SIZE
R05316            INTO  YPDCPGPF-REF-PERIOD-DATE
R05316       END-STRING
R05316       IF   PGPF-REFER-PERIOD-DATE-200 not numeric
R05316       or   PGPF-REFER-PERIOD-DATE-200 = ZEROES
R05316           MOVE '0001-01-01'      TO YPDCPGPF-REF-PERIOD-DATE
R05316       end-if
R05316     END-IF
           MOVE '7601'               TO YPDCPGPF-FINANCIAL-ISTIT
           STRING PGPF-VALUE-DATE(1:4) '-'
                  PGPF-VALUE-DATE(5:2) '-'
                  PGPF-VALUE-DATE(7:2)
                  DELIMITED BY SIZE
                  INTO YPDCPGPF-DATE-VALUE
           END-STRING
      *
DBG==>*    DISPLAY 'DATA VALUE ' YPDCPGPF-DATE-VALUE
      *
            EXEC SQL
              SET :YPDCPGPF-TMSTP-INS = CURRENT TIMESTAMP
            END-EXEC
           .
       F-IMPO-PGPF.
           EXIT.
      *================================================================*
R08421 IMPO-PGPF-301.
R08421*
R08421     INITIALIZE DCLYPTBPGPF.
R08421     STRING COM-DATE-TIME-H(1:4) '-'
R08421            COM-DATE-TIME-H(5:2) '-'
R08421            COM-DATE-TIME-H(7:2)
R08421            DELIMITED BY SIZE
R08421            INTO  YPDCPGPF-PGPF-DATE
R08421     END-STRING.
R08421     STRING COM-DATE-TIME-H(9:2)  ':'
R08421            COM-DATE-TIME-H(11:2) ':'
R08421            COM-DATE-TIME-H(13:2)
R08421            DELIMITED BY SIZE
R08421            INTO  YPDCPGPF-PGPF-TIME
R08421     END-STRING.
      *
      *COLONNE DELLA TABELLA NON VALORIZZATE PER FUNCTION CODE 301
      *
R08421*    MOVE                      TO YPDCPGPF-SWFT-CODE
R08421*    MOVE                      TO YPDCPGPF-PAYMENT-NUMBER
R08421*    MOVE                      TO YPDCPGPF-PAYMENT-SCHEME
R08421*    MOVE                      TO YPDCPGPF-VAT-NUM
R08421* La move seguente � stata fatta per function code 200 e 300
R08421*    MOVE                      TO YPDCPGPF-REF-PERIOD-DATE
R08421* La move seguente � stata fatta per function code 200
R08421*    MOVE                      TO YPDCPGPF-DEVICE-TOT
R08421* La move seguente � stata fatta per function code 200
R08421*    MOVE                      TO YPDCPGPF-UNIT-AMOUNT
R08421     MOVE PGPF-SUMMARY-UID-301 TO YPDCPGPF-SUMMARY-UID
DBG==>*    DISPLAY 'PGPF-FUNCT-CODE-301(' YPDCPGPF-FUNCTION-CODE')'
R08421     MOVE PGPF-FUNCT-CODE-301  TO YPDCPGPF-FUNCTION-CODE
R08421     MOVE PGPF-PAYEMT-UID-301  TO YPDCPGPF-PAYMENT-UID
R08421     MOVE PGPF-SRC-COD-IND-301 TO YPDCPGPF-SOURCE-CODE-IND
R08421     MOVE PGPF-BRAND-CODE-301  TO YPDCPGPF-BRAND
R08421     MOVE PGPF-BANK-ACC-TYP-301
R08421       TO YPDCPGPF-BANK-ACCOUNT-TYPE
R08421     MOVE PGPF-PAYMT-TYPE-301  TO YPDCPGPF-PAYMENT-TYPE
R08421     MOVE PGPF-PAYMT-CYCL-301  TO YPDCPGPF-PAYMENT-CYCLE
R08421     MOVE PGPF-DB-CR-FLAG-301  TO YPDCPGPF-D-C-IND
R08421     COMPUTE YPDCPGPF-PAYMENT-AMOUNT =
R08421             PGPF-PAYMT-TOT-301 / 100
R08421     MOVE PGPF-BILL-CURR-CODE-301
R08421       TO YPDCPGPF-BILLING-CURRENCY
R08421     MOVE PGPF-BANK-ACCOUNT-301
R08421       TO YPDCPGPF-BANK-ACCOUNT-NUM
R08421     MOVE PGPF-LEVEL-PAY-CODE-301  TO YPDCPGPF-LEVEL-PAYMENT-CODE
R08421     MOVE PGPF-ME-ID-CODE-301      TO YPDCPGPF-MERCHANT-ID
R08421     MOVE PGPF-NAME-301            TO YPDCPGPF-NOME-MERCHANT
R08421     MOVE PGPF-BILLING-FLAG-301    TO YPDCPGPF-NET-GROSS-IND
R08421     MOVE PGPF-SUPPRESS-CODE-301   TO YPDCPGPF-SUPPRESS-CODE
R03323     MOVE PGPF-BATCH-ID-301        TO YPDCPGPF-BATCH-ID
R08421     MOVE '7601'               TO YPDCPGPF-FINANCIAL-ISTIT
R08421     STRING PGPF-VALUE-DATE(1:4) '-'
R08421            PGPF-VALUE-DATE(5:2) '-'
R08421            PGPF-VALUE-DATE(7:2)
R08421            DELIMITED BY SIZE
R08421            INTO YPDCPGPF-DATE-VALUE
R08421     END-STRING
R08421* INSERITO PER DEBUG MA BISOGNA CAPIRE E IMPOSTARE IL VALORE
R08421* CORRETTO
R08421     MOVE '0001-01-01'         TO YPDCPGPF-REF-PERIOD-DATE
R08421*
R08421      EXEC SQL
R08421        SET :YPDCPGPF-TMSTP-INS = CURRENT TIMESTAMP
R08421      END-EXEC
R08421     .
R08421 F-IMPO-PGPF-301.
R08421     EXIT.
      *================================================================*
       INSE-PGPF.
      *
            EXEC SQL INSERT INTO YPTBPGPF
            (
             TMSTP_INS
            ,PGPF_DATE
            ,PGPF_TIME
            ,FUNCTION_CODE
            ,PAYMENT_UID
            ,SOURCE_CODE_IND
            ,BRAND
            ,BANK_ACCOUNT_TYPE
            ,PAYMENT_TYPE
            ,PAYMENT_CYCLE
            ,D_C_IND
            ,PAYMENT_AMOUNT
            ,BILLING_CURRENCY
            ,BANK_ACCOUNT_NUM
            ,SWFT_CODE
            ,PAYMENT_NUMBER
            ,LEVEL_PAYMENT_CODE
            ,MERCHANT_ID
            ,NOME_MERCHANT
            ,NET_GROSS_IND
            ,SUPPRESS_CODE
            ,PAYMENT_SCHEME
            ,VAT_NUM
            ,FINANCIAL_ISTIT
            ,DATE_VALUE
R05316      ,REF_PERIOD_DATE
R05316      ,DEVICE_TOT
R05316      ,UNIT_AMOUNT
R08421      ,SUMMARY_UID
R03323      ,BATCH_ID
R54824      ,PAYMENT_REASON
             )
            VALUES (
              :YPDCPGPF-TMSTP-INS
             ,:YPDCPGPF-PGPF-DATE
             ,:YPDCPGPF-PGPF-TIME
             ,:YPDCPGPF-FUNCTION-CODE
             ,:YPDCPGPF-PAYMENT-UID
             ,:YPDCPGPF-SOURCE-CODE-IND
             ,:YPDCPGPF-BRAND
             ,:YPDCPGPF-BANK-ACCOUNT-TYPE
             ,:YPDCPGPF-PAYMENT-TYPE
             ,:YPDCPGPF-PAYMENT-CYCLE
             ,:YPDCPGPF-D-C-IND
             ,:YPDCPGPF-PAYMENT-AMOUNT
             ,:YPDCPGPF-BILLING-CURRENCY
             ,:YPDCPGPF-BANK-ACCOUNT-NUM
             ,:YPDCPGPF-SWFT-CODE
             ,:YPDCPGPF-PAYMENT-NUMBER
             ,:YPDCPGPF-LEVEL-PAYMENT-CODE
             ,:YPDCPGPF-MERCHANT-ID
             ,:YPDCPGPF-NOME-MERCHANT
             ,:YPDCPGPF-NET-GROSS-IND
             ,:YPDCPGPF-SUPPRESS-CODE
             ,:YPDCPGPF-PAYMENT-SCHEME
             ,:YPDCPGPF-VAT-NUM
             ,:YPDCPGPF-FINANCIAL-ISTIT
             ,:YPDCPGPF-DATE-VALUE
R05316       ,:YPDCPGPF-REF-PERIOD-DATE
R05316       ,:YPDCPGPF-DEVICE-TOT
R05316       ,:YPDCPGPF-UNIT-AMOUNT
R08421       ,:YPDCPGPF-SUMMARY-UID
R03323       ,:YPDCPGPF-BATCH-ID
R54824       ,:YPDCPGPF-PAYMENT-REASON
                                )
           END-EXEC.
      *
           IF SQLCODE = 0
              ADD 1                         TO CTR-TABPGPF-INSE
           ELSE
              MOVE SQLCODE                  TO W100-APPO-SQLCODE
              IF SQLCODE = -803
                 ADD 1                      TO W300-SCART-DUPKEY
                 SET WS-SCRI-SCAR-SI        TO TRUE
                 SET WS-SALT-ELAB-SI        TO TRUE
                 PERFORM   IMPO-ERRO-X-DUP-KEY
                    THRU F-IMPO-ERRO-X-DUP-KEY
              ELSE
                 PERFORM   IMPO-ERRO-INSE
                    THRU F-IMPO-ERRO-INSE
              END-IF
           END-IF.
      *
       F-INSE-PGPF.
           EXIT.
      *================================================================*
       CTRL-FINA.
      *
      *--* Controlla se il file � vuoto
           IF ULT-REC-CODA OR ULT-REC-TESTA-680
              CONTINUE
           ELSE
              PERFORM   IMPO-ERRO-ULTI-RECO
                 THRU F-IMPO-ERRO-ULTI-RECO
           END-IF
           .
       F-CTRL-FINA.
           EXIT.
      *================================================================*
      *    Routine di segnalazione di fine programma                   *
      *    Scrive record di stampa di record letti/scritti             *
      *================================================================*
       STAM-RIGH-TOTA.
      *
           MOVE SPACES                      TO YPCWS001-RIGA
           PERFORM SCRIVI-ST              THRU EX-SCRIVI-ST
      *
           PERFORM IMPO-RIGH-TOTA         THRU F-IMPO-RIGH-TOTA
      *
           PERFORM STMP-RIGH-T1           THRU F-STMP-RIGH-T1
      *
R05316*    IF   CTR-CONT-LETTI-SCAR-APO = 0 AND
R05316*         CTR-CONT-LETTI-SCAR-GPO = 0 AND
R05316*         CTR-CONT-LETTI-SCAR-MPO = 0 AND
R05316*         CTR-CONT-LETTI-SCAR-POS = 0 AND
R05316*         CTR-CONT-LETTI-SCAR-VPO = 0
R05316     IF   CTR-CONT-LETTI-SCAR-BILL = 0
                CONTINUE
           ELSE
              PERFORM STMP-RIGH-T2        THRU F-STMP-RIGH-T2
           END-IF
      *
           PERFORM STMP-RIGH-T3           THRU F-STMP-RIGH-T3
      *
           PERFORM STMP-RIGH-T5           THRU F-STMP-RIGH-T5
      *
R14316     PERFORM STMP-RIGH-T6           THRU F-STMP-RIGH-T6
R14316     PERFORM STMP-RIGH-T7           THRU F-STMP-RIGH-T7
R11817     PERFORM STMP-RIGH-T8           THRU F-STMP-RIGH-T8
R05818     PERFORM STMP-RIGH-T9           THRU F-STMP-RIGH-T9
R10219     PERFORM STMP-RIGH-T10          THRU F-STMP-RIGH-T10
R11422     PERFORM STMP-RIGH-T11          THRU F-STMP-RIGH-T11
TK1274     PERFORM STMP-RIGH-T13          THRU F-STMP-RIGH-T13
R11422     PERFORM STMP-RIGH-T12          THRU F-STMP-RIGH-T12
      *
           PERFORM STMP-RIGH-T4           THRU F-STMP-RIGH-T4
           .
       F-STAM-RIGH-TOTA.
           EXIT.
      *================================================================*
       IMPO-RIGH-TOTA.
      *
           MOVE CTR-CONT-LETTI-HEAD-680     TO ETR-CONT-LETTI-HEAD-680
           MOVE CTR-CONT-LETTI-HEAD-681     TO ETR-CONT-LETTI-HEAD-681
           MOVE CTR-CONT-LETTI-TRAIL        TO ETR-CONT-LETTI-TRAIL
           MOVE CTR-CONT-LETTI-TOT          TO ETR-CONT-LETTI-TOT
           MOVE CTR-CONT-LETTI-SCAR         TO ETR-CONT-LETTI-SCAR
           MOVE CTR-CONT-LETTI-DATI-200     TO ETR-CONT-LETTI-DATI-200
           MOVE CTR-CONT-LETTI-DATI-300     TO ETR-CONT-LETTI-DATI-300
R08421     MOVE CTR-CONT-LETTI-DATI-301     TO ETR-CONT-LETTI-DATI-301
           MOVE CTR-CONT-SCARTI             TO ETR-CONT-SCARTI
R05818     MOVE CTR-CONT-FXML               TO ETR-CONT-FXML
R11422     MOVE CTR-CONT-FXM2               TO ETR-CONT-FXM2
           MOVE CTR-CONT-SCRITTI            TO ETR-CONT-SCRITTI
           MOVE CTR-CONT-SCRITTI-B          TO ETR-CONT-SCRITTI-B
R14316     MOVE CTR-CONT-SCRITTI-CONT       TO ETR-CONT-SCRITTI-CONT
R14316     MOVE CTR-CONT-SCRITTI-DCD        TO ETR-CONT-SCRITTI-DCD
R12019     MOVE CTR-CONT-SCRITTI-BILLCCB    TO ETR-CONT-SCRITTI-BILLCCB
R05316*    MOVE CTR-CONT-LETTI-SCAR-APO     TO ETR-CONT-LETTI-SCAR-APO
R05316*    MOVE CTR-CONT-LETTI-SCAR-GPO     TO ETR-CONT-LETTI-SCAR-GPO
R05316*    MOVE CTR-CONT-LETTI-SCAR-MPO     TO ETR-CONT-LETTI-SCAR-MPO
R05316*    MOVE CTR-CONT-LETTI-SCAR-POS     TO ETR-CONT-LETTI-SCAR-POS
R05316*    MOVE CTR-CONT-LETTI-SCAR-VPO     TO ETR-CONT-LETTI-SCAR-VPO
           MOVE CTR-CONT-LETTI-SCAR-BILL    TO ETR-CONT-LETTI-SCAR-BILL
           MOVE CTR-TABPGPF-INSE            TO ETR-TABPGPF-INSE
R11422     MOVE CTR-TABFAS2-INSE            TO ETR-TABFAS2-INSE
TK1274     MOVE CTR-TABFPR-NOT-FOUND        TO ETR-TABFPR-NOT-FOUND
TK1274     MOVE CTR-TABFPR-LETTE            TO ETR-TABFPR-LETTE
           MOVE CTR-CONT-SCRITTI-YPOERRO    TO ETR-CONT-SCRITTI-YPOERRO
           .
       F-IMPO-RIGH-TOTA.
           EXIT.
      *================================================================*
       STMP-RIGH-T1.
      *
           MOVE    '**** RIEPILOGO OPERAZIONI DI INPUT OUTPUT ****'
                                   TO   YPCWS001-RIGA
           PERFORM SCRIVI-ST  THRU      EX-SCRIVI-ST
      *
           MOVE    SPACES     TO        YPCWS001-RIGA
           PERFORM SCRIVI-ST  THRU      EX-SCRIVI-ST
      *
           MOVE    SPACES     TO        YPCWS001-RIGA
           STRING  'TOTALE REC. TESTA LETTI (A.C.680)__:'
                   ETR-CONT-LETTI-HEAD-680
                   DELIMITED BY SIZE
                   INTO YPCWS001-RIGA
           END-STRING
           PERFORM SCRIVI-ST  THRU      EX-SCRIVI-ST
      *
           MOVE    SPACES     TO        YPCWS001-RIGA
           STRING  'TOTALE REC. TESTA LETTI (A.C.681)__:'
                   ETR-CONT-LETTI-HEAD-681
                   DELIMITED BY SIZE
                   INTO YPCWS001-RIGA
           END-STRING
           PERFORM SCRIVI-ST  THRU      EX-SCRIVI-ST
      *
           MOVE    SPACES     TO        YPCWS001-RIGA
           STRING  'TOTALE RECORDS DI CODA LETTI_______:'
                   ETR-CONT-LETTI-TRAIL
                   DELIMITED BY SIZE
                   INTO YPCWS001-RIGA
           END-STRING
           PERFORM SCRIVI-ST  THRU      EX-SCRIVI-ST
      *
           MOVE    SPACES     TO        YPCWS001-RIGA
           STRING  'TOTALE RECORDS DATI (F.C.200) LETTI:'
                   ETR-CONT-LETTI-DATI-200
                   DELIMITED BY SIZE
                   INTO YPCWS001-RIGA
           END-STRING
           PERFORM SCRIVI-ST  THRU      EX-SCRIVI-ST
      *
           MOVE    SPACES     TO        YPCWS001-RIGA
           STRING  'TOTALE RECORDS DATI (F.C.300) LETTI:'
                   ETR-CONT-LETTI-DATI-300
                   DELIMITED BY SIZE
                   INTO YPCWS001-RIGA
           END-STRING
           PERFORM SCRIVI-ST  THRU      EX-SCRIVI-ST
      *
R08421     MOVE    SPACES     TO        YPCWS001-RIGA
R08421     STRING  'TOTALE RECORDS DATI (F.C.301) LETTI:'
R08421             ETR-CONT-LETTI-DATI-301
R08421             DELIMITED BY SIZE
R08421             INTO YPCWS001-RIGA
R08421     END-STRING
R08421     PERFORM SCRIVI-ST  THRU      EX-SCRIVI-ST
      *
           MOVE    SPACES     TO        YPCWS001-RIGA
           STRING  'TOTALE REC.DATI SCARTATI IN LETTURA:'
                   ETR-CONT-LETTI-SCAR
                   DELIMITED BY SIZE
                   INTO YPCWS001-RIGA
           END-STRING
           PERFORM SCRIVI-ST  THRU      EX-SCRIVI-ST
      *
           .
       F-STMP-RIGH-T1.
           EXIT.
      *================================================================*
       STMP-RIGH-T2.
      *
           MOVE    SPACES     TO        YPCWS001-RIGA
R05316     IF CTR-CONT-LETTI-SCAR-BILL NOT = 0
R05316        MOVE ' CANONI BILL_______________:' TO RIGA-DI-CUI-RESTO
R05316        MOVE  CTR-CONT-LETTI-SCAR-BILL      TO RIGA-DI-CUI-NUM
R05316        COMPUTE CTR-COM-IMPORTO = CTR-CONT-LETTI-SCAR-BILLI / 100
R05316        MOVE  CTR-COM-IMPORTO         TO RIGA-DI-CUI-IMP
R05316        MOVE  RIGA-DI-CUI             TO YPCWS001-RIGA
R05316        PERFORM SCRIVI-ST  THRU      EX-SCRIVI-ST
R05316        MOVE    SPACES     TO        YPCWS001-RIGA
R05316        MOVE    '       .' TO        FILLER-YY
R05316     END-IF
      *
R05316*    IF CTR-CONT-LETTI-SCAR-APO NOT = 0
R05316*       MOVE ' CANONI APO________________:' TO RIGA-DI-CUI-RESTO
R05316*       MOVE  CTR-CONT-LETTI-SCAR-APO       TO RIGA-DI-CUI-NUM
R05316*       COMPUTE CTR-COM-IMPORTO = CTR-CONT-LETTI-SCAR-APOI / 100
R05316*       MOVE  CTR-COM-IMPORTO         TO RIGA-DI-CUI-IMP
R05316*       MOVE  RIGA-DI-CUI             TO YPCWS001-RIGA
R05316*       PERFORM SCRIVI-ST  THRU      EX-SCRIVI-ST
R05316*       MOVE    SPACES     TO        YPCWS001-RIGA
R05316*       MOVE    '       .' TO        FILLER-YY
R05316*    END-IF
      *
R05316*    IF CTR-CONT-LETTI-SCAR-GPO NOT = 0
R05316*       MOVE ' CANONI GPO________________:' TO RIGA-DI-CUI-RESTO
R05316*       MOVE  CTR-CONT-LETTI-SCAR-GPO       TO RIGA-DI-CUI-NUM
R05316*       COMPUTE CTR-COM-IMPORTO = CTR-CONT-LETTI-SCAR-GPOI / 100
R05316*       MOVE  CTR-COM-IMPORTO         TO RIGA-DI-CUI-IMP
R05316*       MOVE  RIGA-DI-CUI             TO YPCWS001-RIGA
R05316*       PERFORM SCRIVI-ST  THRU      EX-SCRIVI-ST
R05316*       MOVE    SPACES     TO        YPCWS001-RIGA
R05316*       MOVE    '       .' TO        FILLER-YY
R05316*    END-IF
      *
R05316*    IF CTR-CONT-LETTI-SCAR-MPO NOT = 0
R05316*       MOVE ' CANONI MPO________________:' TO RIGA-DI-CUI-RESTO
R05316*       MOVE  CTR-CONT-LETTI-SCAR-MPO       TO RIGA-DI-CUI-NUM
R05316*       COMPUTE CTR-COM-IMPORTO = CTR-CONT-LETTI-SCAR-MPOI / 100
R05316*       MOVE  CTR-COM-IMPORTO         TO RIGA-DI-CUI-IMP
R05316*       MOVE  RIGA-DI-CUI             TO YPCWS001-RIGA
R05316*       PERFORM SCRIVI-ST  THRU      EX-SCRIVI-ST
R05316*       MOVE    SPACES     TO        YPCWS001-RIGA
R05316*       MOVE    '       .' TO        FILLER-YY
R05316*    END-IF
      *
R05316*    IF CTR-CONT-LETTI-SCAR-POS NOT = 0
R05316*       MOVE ' CANONI POS________________:' TO RIGA-DI-CUI-RESTO
R05316*       MOVE  CTR-CONT-LETTI-SCAR-POS       TO RIGA-DI-CUI-NUM
R05316*       COMPUTE CTR-COM-IMPORTO = CTR-CONT-LETTI-SCAR-POSI / 100
R05316*       MOVE  CTR-COM-IMPORTO         TO RIGA-DI-CUI-IMP
R05316*       MOVE  RIGA-DI-CUI             TO YPCWS001-RIGA
R05316*       PERFORM SCRIVI-ST  THRU      EX-SCRIVI-ST
R05316*       MOVE    SPACES     TO        YPCWS001-RIGA
R05316*       MOVE    '       .' TO        FILLER-YY
R05316*    END-IF
R05316*
R05316*    IF CTR-CONT-LETTI-SCAR-VPO NOT = 0
R05316*       MOVE ' CANONI VPO________________:' TO RIGA-DI-CUI-RESTO
R05316*       MOVE  CTR-CONT-LETTI-SCAR-VPO       TO RIGA-DI-CUI-NUM
R05316*       COMPUTE CTR-COM-IMPORTO = CTR-CONT-LETTI-SCAR-VPOI / 100
R05316*       MOVE  CTR-COM-IMPORTO         TO RIGA-DI-CUI-IMP
R05316*       MOVE  RIGA-DI-CUI             TO YPCWS001-RIGA
R05316*       PERFORM SCRIVI-ST  THRU      EX-SCRIVI-ST
R05316*       MOVE    SPACES     TO        YPCWS001-RIGA
R05316*       MOVE    '       .' TO        FILLER-YY
R05316*    END-IF
           .
       F-STMP-RIGH-T2.
           EXIT.
      *================================================================*
       STMP-RIGH-T3.
      *
           MOVE    SPACES     TO        YPCWS001-RIGA
           STRING  'TOTALE RECORDS LETTI_______________:'
                   ETR-CONT-LETTI-TOT
                   DELIMITED BY SIZE
                   INTO YPCWS001-RIGA
           END-STRING
           PERFORM SCRIVI-ST  THRU      EX-SCRIVI-ST
      *
           MOVE    SPACES     TO        YPCWS001-RIGA
           STRING  'TOTALE RECORDS D50 ONLINE SCRITTI__:'
                   ETR-CONT-SCRITTI
                   DELIMITED BY SIZE
                   INTO YPCWS001-RIGA
           END-STRING
           PERFORM SCRIVI-ST  THRU      EX-SCRIVI-ST
      *
           MOVE    SPACES     TO        YPCWS001-RIGA
           STRING  'TOTALE RECORDS D50 BATCH SCRITTI___:'
                   ETR-CONT-SCRITTI-B
                   DELIMITED BY SIZE
                   INTO YPCWS001-RIGA
           END-STRING
           PERFORM SCRIVI-ST  THRU      EX-SCRIVI-ST
      *
           MOVE    SPACES     TO        YPCWS001-RIGA
           STRING  'TOTALE RECORDS SU FILE SCARTI______:'
                   ETR-CONT-SCARTI
                   DELIMITED BY SIZE
                   INTO YPCWS001-RIGA
           END-STRING
           PERFORM SCRIVI-ST  THRU      EX-SCRIVI-ST
      *
           .
       F-STMP-RIGH-T3.
           EXIT.
      *================================================================*
       STMP-RIGH-T5.
      *
      *
           MOVE    SPACES     TO        YPCWS001-RIGA
           STRING  'TOTALE RECORDS SCRITTI YPOERRO_____:'
                   ETR-CONT-SCRITTI-YPOERRO
                   DELIMITED BY SIZE
                   INTO YPCWS001-RIGA
           END-STRING
           PERFORM SCRIVI-ST  THRU      EX-SCRIVI-ST
           .
       F-STMP-RIGH-T5.
           EXIT.
      *================================================================*
R14316 STMP-RIGH-T6.
      *
      *
           MOVE    SPACES     TO        YPCWS001-RIGA
           STRING  'TOTALE RECORDS SCRITTI XYDCONT_____:'
                   ETR-CONT-SCRITTI-CONT
                   DELIMITED BY SIZE
                   INTO YPCWS001-RIGA
           END-STRING
           PERFORM SCRIVI-ST  THRU      EX-SCRIVI-ST
           .
       F-STMP-RIGH-T6.
           EXIT.
      *================================================================*
R14316 STMP-RIGH-T7.
      *
      *
           MOVE    SPACES     TO        YPCWS001-RIGA
           STRING  'TOTALE RECORDS SCRITTI OUTDCD______:'
                   ETR-CONT-SCRITTI-DCD
                   DELIMITED BY SIZE
                   INTO YPCWS001-RIGA
           END-STRING
           PERFORM SCRIVI-ST  THRU      EX-SCRIVI-ST
           .
       F-STMP-RIGH-T7.
           EXIT.
      *================================================================*
R11817 STMP-RIGH-T8.
      *
           MOVE    SPACES     TO        YPCWS001-RIGA
           STRING  'TOTALE RECORDS FILE SCARTI Evo.Bus.:'
                   ETR-CONT-SCRITTI-DCD
                   DELIMITED BY SIZE
                   INTO YPCWS001-RIGA
           END-STRING
           PERFORM SCRIVI-ST  THRU      EX-SCRIVI-ST
           .
R11817 F-STMP-RIGH-T8.
           EXIT.
      *================================================================*
R05818 STMP-RIGH-T9.
      *
           MOVE    SPACES     TO        YPCWS001-RIGA
           STRING  'TOTALE RECORDS FILE FXML:           '
                   ETR-CONT-FXML
                   DELIMITED BY SIZE
                   INTO YPCWS001-RIGA
           END-STRING
           PERFORM SCRIVI-ST  THRU      EX-SCRIVI-ST
           .
R05818 F-STMP-RIGH-T9.
           EXIT.
      *================================================================*
R12019 STMP-RIGH-T10.
      *
           MOVE    SPACES     TO        YPCWS001-RIGA
           STRING  'TOTALE RECORDS SCRITTI BILLCCB_____ '
                   ETR-CONT-SCRITTI-BILLCCB
                   DELIMITED BY SIZE
                   INTO YPCWS001-RIGA
           END-STRING
           PERFORM SCRIVI-ST  THRU      EX-SCRIVI-ST
           .
R12019 F-STMP-RIGH-T10.
           EXIT.
      *================================================================*
       STMP-RIGH-T4.
      *
           MOVE    SPACES     TO        YPCWS001-RIGA
           PERFORM SCRIVI-ST  THRU      EX-SCRIVI-ST
      *
           MOVE    SPACES     TO        YPCWS001-RIGA
           STRING  'RIGHE INSERITE IN TABELLA YPTBPGPF_:'
                  ETR-TABPGPF-INSE
                  DELIMITED BY SIZE
                  INTO  YPCWS001-RIGA
           END-STRING.
           PERFORM SCRIVI-ST THRU EX-SCRIVI-ST
           .
       F-STMP-RIGH-T4.
           EXIT.
      *================================================================*
TK1274 STMP-RIGH-T13.
      *
           MOVE    SPACES     TO        YPCWS001-RIGA
           STRING  'RIGHE NON TROVATE IN GEP FPR ______:'
                  ETR-TABFPR-NOT-FOUND
                  DELIMITED BY SIZE
                  INTO  YPCWS001-RIGA
           END-STRING.
           PERFORM SCRIVI-ST THRU EX-SCRIVI-ST
           MOVE    SPACES     TO        YPCWS001-RIGA
           STRING  'RIGHE TROVATE IN GEP FPR __________:'
                  ETR-TABFPR-LETTE
                  DELIMITED BY SIZE
                  INTO  YPCWS001-RIGA
           END-STRING.
           PERFORM SCRIVI-ST THRU EX-SCRIVI-ST
           .
TK1274 F-STMP-RIGH-T13.
           EXIT.
      *================================================================*
       STAM-RIGH-FINA.
      *
           MOVE SPACES                      TO YPCWS001-RIGA
           MOVE SPACES                      TO YP-MSGERR
      *
           PERFORM IMPO-ORA               THRU EX-IMPO-ORA
      *
           STRING 'ELABORAZIONE CHIUSA CORRETTAMENTE - DATA: '
                  WS-DATA-GGMMSSAA
                  '  ORA: ' WS-ORA-DAY-2
                  DELIMITED BY SIZE
                  INTO  YP-MSGERR-2
           END-STRING.
      *
           MOVE 3                           TO YP-INDMAX
      *
           PERFORM YP-SCRIVIERR           THRU EX-YP-SCRIVIERR
           .
       F-STAM-RIGH-FINA.
           EXIT.
      *================================================================
      * CHIUSURA ARCHIVI
      *================================================================
       CLOSE-FILE.
      *
           CLOSE IPAYMENT
FIANNH           OPECONTB
                 OPECONT
R14316           XYDCONT
R14316           OUTDCD
                 OUSCARTI
R05818           OUTFXML
R11422           OUTFXM2
                 YYDTABE
                 YPOERRO
                 ST
R12019           BILLCCB
                 .
       F-CLOSE-FILE.
           EXIT.
      *==============================================================*
       IMPO-ORA.
      *
           MOVE   SPACES                    TO WS-ORA-DAY-2
      *
      *--* Acquisizione ora
           ACCEPT WS-ORA-DAY              FROM TIME
           STRING WS-ORA-DAY(1:2) ':'
                  WS-ORA-DAY(3:2) ':'
                  WS-ORA-DAY(5:2)
           DELIMITED BY SIZE
             INTO WS-ORA-DAY-2
           END-STRING
           .
       EX-IMPO-ORA.
           EXIT.
      *==============================================================*
      *
      *==============================================================*
       YP-SCRIVIERR.
      *
      *--* Scarica tabella messaggi su stampa
           PERFORM VARYING YP-IND FROM 1 BY 1
                   UNTIL   YP-IND > YP-INDMAX
             MOVE YP-MSGOCC(YP-IND)   TO     YPCWS001-RIGA
             PERFORM SCRIVI-ST THRU EX-SCRIVI-ST
           END-PERFORM.
       EX-YP-SCRIVIERR.
           EXIT.
      *================================================================*
       IMPO-ERRO-INCC-HIVA.
      *
           MOVE SPACES                      TO YPCWS001-RIGA
           MOVE SPACES                      TO YP-MSGERR
      *
           MOVE 'ERRORE GRAVE CALL PGM. CRVYD228'
                                            TO P303-MSGER-DESCR
           MOVE INCC-RETCODE                TO P303-FILE-STATUS
           MOVE INCC-CV20-RAPPORT           TO P303-MSGER-DATO
           PERFORM ERRORE-P303            THRU EX-ERRORE-P303
           .
       F-IMPO-ERRO-INCC-HIVA.
           EXIT.
      *================================================================*
       IMPO-ERRO-INCC-NO.
      *
           MOVE SPACES                      TO YPCWS001-RIGA
           MOVE SPACES                      TO YP-MSGERR
      *
           STRING 'RAPPORTO: ' INCC-CV20-RAPPORT
                  DELIMITED BY SIZE
                  INTO YPCWS001-RIGA
           END-STRING
      *
           PERFORM SCRIVI-ST              THRU EX-SCRIVI-ST
           MOVE 'ERRORE CRVYD228 '          TO P303-MSGER-DESCR
           MOVE INCC-RETCODE                TO P303-FILE-STATUS
           MOVE INCC-CV20-RAPPORT           TO P303-MSGER-DATO
           PERFORM ERRORE-P303            THRU EX-ERRORE-P303
           .
       F-IMPO-ERRO-INCC-NO.
           EXIT.
      *================================================================*
       IMPO-ERRO-INCC-NF.
      *
           MOVE SPACES                      TO YPCWS001-RIGA
           STRING 'IBAN: ' PGPF-BANK-ACCOUNT ' - '
                  'RAPPORTO: ' INCC-CV20-RAPPORT
                  DELIMITED BY SIZE
                  INTO YPCWS001-RIGA
           END-STRING
           PERFORM SCRIVI-ST              THRU EX-SCRIVI-ST
      *
           MOVE SPACES                      TO YPCWS001-RIGA
           STRING 'FRAZIONARIO NON TROVATO - '
                  'RECORD SCARTATO'
                  DELIMITED BY SIZE
                  INTO YPCWS001-RIGA
           END-STRING
           PERFORM SCRIVI-ST              THRU EX-SCRIVI-ST
           .
       F-IMPO-ERRO-INCC-NF.
           EXIT.
      *================================================================*
       IMPO-ERRO-FCOD-TEST.
      *
           MOVE SPACES                      TO YPCWS001-RIGA
           MOVE SPACES                      TO YP-MSGERR
      *
           MOVE 'IL PRIMO RECORD NON E'' UN RECORD DI TESTA '
                                            TO YP-MSGERR-1
      *
           STRING 'FUNCTION CODE DEL PRIMO RECORD E'' UGUALE A '
                  PGPFH-FUNCT-CODE
           DELIMITED BY SIZE              INTO YP-MSGERR-2
           END-STRING
      *
           PERFORM GEST-ERRO-SU-TRE-RIGH  THRU F-GEST-ERRO-SU-TRE-RIGH
           .
       F-IMPO-ERRO-FCOD-TEST.
           EXIT.
      *================================================================*
       IMPO-ERRO-CODA.
      *
           MOVE SPACES                      TO YPCWS001-RIGA
           MOVE SPACES                      TO YP-MSGERR
      *
           MOVE 'RECORD DI CODA SENZA RECORD DI TESTA'
                                            TO YP-MSGERR-1
      *
           STRING 'ID. RECORD  DI CODA ==> ' YPCRPGPF-TRAILER(1:23)
                  DELIMITED BY SIZE INTO YP-MSGERR-2
           END-STRING
      *
           PERFORM IMPO-ORA          THRU  EX-IMPO-ORA
           STRING 'ELABORAZIONE INTERROTTA  - DATA: '
                  WS-DATA-GGMMSSAA
                  '  ORA: ' WS-ORA-DAY-2
                  DELIMITED BY SIZE
                  INTO  YP-MSGERR-4
           END-STRING
           PERFORM ROLL-BACK              THRU F-ROLL-BACK
           PERFORM IMPO-RIG4              THRU F-IMPO-RIG4
           .
       F-IMPO-ERRO-CODA.
           EXIT.
      *================================================================*
       IMPO-ERRO-TEST.
      *
           MOVE SPACES                      TO YPCWS001-RIGA
           MOVE SPACES                      TO YP-MSGERR
      *
           MOVE 'FLUSSO PGPF ERRATO VERIFICARE I RECORD DI TESTA E CODA'
                                            TO YP-MSGERR-1
      *
           STRING 'ID. RECORD  MSG. NUMBER: ' PGPFH-MSG-NUMBER
                  '  -DATA: '          PGPFH-DATE-TIME-CRE
                  DELIMITED BY SIZE INTO YP-MSGERR-2
           END-STRING
      *
           PERFORM IMPO-ORA          THRU  EX-IMPO-ORA
           STRING 'ELABORAZIONE INTERROTTA  - DATA: '
                  WS-DATA-GGMMSSAA
                  '  ORA: ' WS-ORA-DAY-2
                  DELIMITED BY SIZE
                  INTO  YP-MSGERR-4
           END-STRING
           PERFORM ROLL-BACK              THRU F-ROLL-BACK
           PERFORM IMPO-RIG4              THRU F-IMPO-RIG4
           .
       F-IMPO-ERRO-TEST.
           EXIT.
      *================================================================*
       IMPO-ERRO-ROUT-TEST.
      *
           MOVE SPACES                      TO YPCWS001-RIGA
           MOVE SPACES                      TO YP-MSGERR
      *
           STRING 'DATA: '
                  PGPFH-DATE-TIME-CRE
                  ' SU RECORD DI TESTA NON VALIDA'
                  ' ERR. => 'YYCWUTDA-FLAG-ERRORE
                  DELIMITED BY SIZE
                  INTO YP-MSGERR-1
           END-STRING
      *
           PERFORM GEST-ERRO-SU-TRE-RIGH  THRU F-GEST-ERRO-SU-TRE-RIGH
           .
       F-IMPO-ERRO-ROUT-TEST.
           EXIT.
      *================================================================*
       IMPO-ERRO-TEST-ACTI-CODE.
      *
           MOVE SPACES                      TO YPCWS001-RIGA
           MOVE SPACES                      TO YP-MSGERR
      *
           STRING 'RECORD DI TESTA ERRATO X ACTION CODE = '
                   PGPFH-ACTION-CODE
                  DELIMITED BY SIZE
                  INTO YP-MSGERR-1
           END-STRING
      *
           PERFORM GEST-ERRO-SU-TRE-RIGH  THRU F-GEST-ERRO-SU-TRE-RIGH
      *
           .
       F-IMPO-ERRO-TEST-ACTI-CODE.
           EXIT.
      *================================================================*
       IMPO-ERRO-FCOD.
      *
           MOVE SPACES                      TO YPCWS001-RIGA
           MOVE SPACES                      TO YP-MSGERR
      *
           MOVE 'FUNCTION CODE ERRATO '
                                            TO YP-MSGERR-1
      *
           STRING 'VALORE TROVATO: ' PGPFH-FUNCT-CODE
           DELIMITED BY SIZE              INTO YP-MSGERR-2
           END-STRING
      *
           PERFORM GEST-ERRO-SU-TRE-RIGH  THRU F-GEST-ERRO-SU-TRE-RIGH
           .
       F-IMPO-ERRO-FCOD.
           EXIT.
      *================================================================*
       IMPO-ERRO-ULTI-RECO.
      *
           MOVE SPACES                      TO YPCWS001-RIGA
           MOVE SPACES                      TO YP-MSGERR
      *
           STRING 'MANCA RECORD DI CODA SULL''ULTIMO '
                  'RECORD DEL FLUSSO PGPF '
                  DELIMITED BY SIZE
                  INTO YP-MSGERR-1
           END-STRING
      *
           PERFORM GEST-ERRO-SU-TRE-RIGH  THRU F-GEST-ERRO-SU-TRE-RIGH
           .
       F-IMPO-ERRO-ULTI-RECO.
           EXIT.
      *================================================================*
       IMPO-ERRO-PRIM-RECO.
      *
           MOVE SPACES                      TO YPCWS001-RIGA
           MOVE SPACES                      TO YP-MSGERR
      *
           MOVE  'FLUSSO - PGPF -      '    TO YP-MSGERR-1
      *
           STRING 'PRIMO RECORD:COD.MSG = ' PGPFH-MSG-TYPE-ID
                  'F.CODE = ' PGPFH-FUNCT-CODE
                  ' NON VALIDO'
                  DELIMITED BY SIZE       INTO YP-MSGERR-2
           END-STRING
      *
           PERFORM IMPO-ORA       THRU  EX-IMPO-ORA
           STRING 'ELABORAZIONE INTERROTTA  - DATA: '
                  WS-DATA-GGMMSSAA
                  '  ORA: ' WS-ORA-DAY-2
                  DELIMITED BY SIZE       INTO YP-MSGERR-4
           END-STRING
           PERFORM IMPO-RIG4              THRU F-IMPO-RIG4
           .
       F-IMPO-ERRO-PRIM-RECO.
           EXIT.
      *================================================================*
       IMPO-ERRO-FILE-VUOTO.
      *
           MOVE SPACES                      TO YPCWS001-RIGA
           MOVE SPACES                      TO YP-MSGERR
      *
           MOVE  'FLUSSO - PGPF -  VUOTO'   TO YP-MSGERR-1
      *
           PERFORM IMPO-ORA               THRU EX-IMPO-ORA
      *
           STRING 'NESSUN DATO ELABORATO - DATA: '
                     WS-DATA-GGMMSSAA
                     '  ORA: ' WS-ORA-DAY-2
           DELIMITED BY SIZE              INTO  YP-MSGERR-3
           END-STRING
      *
           MOVE   3                         TO YP-INDMAX
      *
           PERFORM YP-SCRIVIERR           THRU EX-YP-SCRIVIERR
      *
           GO TO FINE-JOB
           .
       F-IMPO-ERRO-FILE-VUOTO.
           EXIT.
      *================================================================*
       IMPO-ERRO-X-DUP-KEY.
      *
           MOVE SPACES                      TO YPCWS001-RIGA
           STRING 'KEY DOPPIA SU TABELLA DB2 YPTBPGPF-'
                  DELIMITED BY SIZE
                  INTO YPCWS001-RIGA
           END-STRING
           PERFORM SCRIVI-ST              THRU EX-SCRIVI-ST
      *
           MOVE SPACES                      TO YPCWS001-RIGA
           STRING 'PGPF-PAYEMT-UID ==> '   PGPF-PAYEMT-UID
                  DELIMITED BY SIZE
                  INTO YPCWS001-RIGA
           END-STRING
           PERFORM SCRIVI-ST              THRU EX-SCRIVI-ST
      *
           MOVE SPACES                      TO YPCWS001-RIGA
           PERFORM SCRIVI-ST              THRU EX-SCRIVI-ST
           .
      *
      *--* Imposta area x messaggio errori via mail
           MOVE 'Key doppia su tabella DB2 YPTBPGPF-'
                                            TO WS-AREA-APPO-YPOE-DESC
           .
      *
       F-IMPO-ERRO-X-DUP-KEY.
           EXIT.
      *================================================================*
R11422 IMPO-ERRO-X-DUP-KEY-FAS2.
      *
           MOVE SPACES                      TO YPCWS001-RIGA
           STRING 'KEY DOPPIA SU TABELLA DB2 YPTBFAS2-'
                  DELIMITED BY SIZE
                  INTO YPCWS001-RIGA
           END-STRING
           PERFORM SCRIVI-ST              THRU EX-SCRIVI-ST
      *
           MOVE SPACES                      TO YPCWS001-RIGA
           STRING 'PAYEMT-UID ==> '   YPDCFAS2-PGPF-PAYMENT-UID
                  DELIMITED BY SIZE
                  INTO YPCWS001-RIGA
           END-STRING
           PERFORM SCRIVI-ST              THRU EX-SCRIVI-ST
           PERFORM SCRIVI-ST              THRU EX-SCRIVI-ST
           STRING 'NUMERO-RAPPORTO ==> '  YPDCFAS2-NUME-RAPP-FA
                  DELIMITED BY SIZE
                  INTO YPCWS001-RIGA
           END-STRING
           PERFORM SCRIVI-ST              THRU EX-SCRIVI-ST
           STRING 'PGPF-BANK-ACC-TYP=> '  YPDCFAS2-PGPF-BANK-ACC-TYP
                  DELIMITED BY SIZE
                  INTO YPCWS001-RIGA
           END-STRING
           PERFORM SCRIVI-ST              THRU EX-SCRIVI-ST
      *
           MOVE SPACES                      TO YPCWS001-RIGA
           PERFORM SCRIVI-ST              THRU EX-SCRIVI-ST
           .
      *
      *--* Imposta area x messaggio errori via mail
           MOVE 'Key doppia su tabella DB2 YPTBFAS2-'
                                            TO WS-AREA-APPO-YPOE-DESC
           .
      *
R11422 F-IMPO-ERRO-X-DUP-KEY-FAS2.
           EXIT.
      *================================================================*
       IMPO-ERRO-INSE.
      *
           MOVE SPACES                      TO YPCWS001-RIGA
           MOVE SPACES                      TO YP-MSGERR
      *
           STRING
                  'ERRORE INSERT IN TABELLA YPTBPGPF SQLCODE: '
                  W100-APPO-SQLCODE ' riga =>'
           DELIMITED BY SIZE              INTO YP-MSGERR-1
           END-STRING
      *
           MOVE DCLYPTBPGPF                 TO YP-MSGERR-2
      *
           PERFORM GEST-ERRO-SU-TRE-RIGH  THRU F-GEST-ERRO-SU-TRE-RIGH
           .
       F-IMPO-ERRO-INSE.
           EXIT.
      *================================================================*
R11422 IMPO-ERRO-INSE-FAS2.
      *
           MOVE SPACES                      TO YPCWS001-RIGA
           MOVE SPACES                      TO YP-MSGERR
      *
           STRING
                  'ERRORE INSERT IN TABELLA YPTBFAS2 SQLCODE: '
                  W100-APPO-SQLCODE ' riga =>'
           DELIMITED BY SIZE              INTO YP-MSGERR-1
           END-STRING
      *
           MOVE YPDCFAS2                    TO YP-MSGERR-2
      *
           PERFORM GEST-ERRO-SU-TRE-RIGH  THRU F-GEST-ERRO-SU-TRE-RIGH
           .
R11422 F-IMPO-ERRO-INSE-FAS2.
           EXIT.
      *================================================================*
      *    Scrittura record errati usato per inviare mail errati       *
      *----------------------------------------------------------------
      * ROUTINE DI GESTIONE ERRORE GENERICO
      *----------------------------------------------------------------
       ERRORE-P303.
      *
           MOVE SPACES                      TO YPCWS001-RIGA
           MOVE P303-MSG-ERRORE-1           TO YPCWS001-RIGA
           PERFORM SCRIVI-ST              THRU EX-SCRIVI-ST
      *
           MOVE SPACES                      TO YPCWS001-RIGA
           MOVE P303-MSG-ERRORE-2           TO YPCWS001-RIGA
           MOVE '*** ELABORAZIONE INTERROTTA ***' TO YPCWS001-RIGA
           PERFORM SCRIVI-ST              THRU EX-SCRIVI-ST
      *
           MOVE SPACES                      TO YPCWS001-RIGA
           PERFORM SCRIVI-ST              THRU EX-SCRIVI-ST
      *
           PERFORM STAM-RIGH-TOTA         THRU F-STAM-RIGH-TOTA
           PERFORM STAM-RIGH-FINA         THRU F-STAM-RIGH-FINA
           MOVE 12 TO RETURN-CODE
           GO TO FINE-JOB
           .
      *
       EX-ERRORE-P303.  EXIT.

      *================================================================*
       GEST-ERRO-SU-TRE-RIGH.
      *
           PERFORM ROLL-BACK              THRU F-ROLL-BACK
           PERFORM IMPO-ELAB-INTE         THRU F-IMPO-ELAB-INTE
           PERFORM IMPO-RIG3              THRU F-IMPO-RIG3
           .
       F-GEST-ERRO-SU-TRE-RIGH.
           EXIT.
      *================================================================*
       ROLL-BACK.
      *
           EXEC SQL ROLLBACK END-EXEC
           .
      *
       F-ROLL-BACK.
           EXIT.
      *================================================================*
       IMPO-ELAB-INTE.
      *
           PERFORM IMPO-ORA THRU EX-IMPO-ORA
      *
           STRING 'ELABORAZIONE INTERROTTA  - DATA: '
                  WS-DATA-GGMMSSAA
                  '  ORA: ' WS-ORA-DAY-2
                  DELIMITED BY SIZE
                  INTO  YP-MSGERR-3
           END-STRING
      *
           .
       F-IMPO-ELAB-INTE.
           EXIT.
      *================================================================*
       IMPO-RIG3.
      *
           MOVE  3                          TO YP-INDMAX
      *
           PERFORM YP-SCRIVIERR           THRU EX-YP-SCRIVIERR
      *
           MOVE 12                          TO RETURN-CODE
      *
           GO TO FINE-JOB
           .
       F-IMPO-RIG3.
           EXIT.
      *================================================================*
       IMPO-RIG4.
      *
           MOVE  4                          TO YP-INDMAX
      *
           PERFORM YP-SCRIVIERR           THRU EX-YP-SCRIVIERR
      *
           MOVE 12                          TO   RETURN-CODE
      *
           GO TO FINE-JOB
           .
       F-IMPO-RIG4.
           EXIT.
      *================================================================*
R05316 SELE-TABE-GEP-PGP.
      *
           INITIALIZE                    HV-TABE
           MOVE 'GEP'                 TO HV-TABE-KNAMTB1
           MOVE 'PGP'                 TO HV-TABE-KVARTB1(1:3)
           IF   PGPF-PAYMT-TYPE = space
                MOVE 'DEF'                 TO HV-TABE-KVARTB1(4:3)
           ELSE
                MOVE PGPF-PAYMT-TYPE       TO HV-TABE-KVARTB1(4:3)
           END-IF

           EXEC SQL SELECT DATI
                INTO :HV-TABE-DATI
                FROM XYTBTABE
           WHERE KNAMTB1 = :HV-TABE-KNAMTB1 AND
                 KVARTB1 = :HV-TABE-KVARTB1
           END-EXEC
           .

      *
           MOVE SQLCODE                  TO W100-APPO-SQLCODE
           EVALUATE SQLCODE
            WHEN +0
                   MOVE HV-TABE-DATI-A      TO YPCRTPGP-DATI
            WHEN OTHER
              MOVE SPACES                      TO YPCWS001-RIGA
              MOVE SPACES                      TO YP-MSGERR
              STRING
                  'ERRORE read TABELLA xytbtabe SQLCODE: '
                  W100-APPO-SQLCODE ' riga =>'
              DELIMITED BY SIZE              INTO YP-MSGERR-1
              END-STRING
              MOVE HV-TABE                     TO YP-MSGERR-2
      *
              PERFORM GEST-ERRO-SU-TRE-RIGH
                 THRU F-GEST-ERRO-SU-TRE-RIGH
           END-EVALUATE.
      *
           .
R05316 F-SELE-TABE-GEP-PGP.
           EXIT.
      *================================================================*
R05818 SELE-TABE-GEP-CCB.
      *
           SET  WS-LETTA-GEP-CCB-SI   TO TRUE
           INITIALIZE                    HV-TABE
           MOVE 'GEP'                 TO HV-TABE-KNAMTB1
           MOVE 'CCB'                 TO HV-TABE-KVARTB1(1:3)

           EXEC SQL SELECT DATI
                INTO :HV-TABE-DATI
                FROM XYTBTABE
           WHERE KNAMTB1 = :HV-TABE-KNAMTB1 AND
                 KVARTB1 = :HV-TABE-KVARTB1
           END-EXEC
           .

      *
           MOVE SQLCODE                  TO W100-APPO-SQLCODE
           EVALUATE SQLCODE
            WHEN +0
                   MOVE HV-TABE-DATI-A      TO YPCRTCCB-DATI
            WHEN OTHER
              MOVE SPACES                      TO YPCWS001-RIGA
              MOVE SPACES                      TO YP-MSGERR
              STRING
                  'ERRORE read TABELLA xytbtabe SQLCODE: '
                  W100-APPO-SQLCODE ' riga =>'
              DELIMITED BY SIZE              INTO YP-MSGERR-1
              END-STRING
              MOVE HV-TABE                     TO YP-MSGERR-2
      *
              PERFORM GEST-ERRO-SU-TRE-RIGH
                 THRU F-GEST-ERRO-SU-TRE-RIGH
           END-EVALUATE.
      *
           .
R05818 F-SELE-TABE-GEP-CCB.
           EXIT.
      *================================================================*
R14217 TROVA-KEY-RANDOM.

           EXEC SQL
R10418*    SELECT SUBSTR(HEX(BIGINT((RAND() * 1000000000))),9,8)
R10418     SELECT BIGINT((RAND() * 1000000000))
            INTO :HV-KEY-RANDOM
            FROM SYSIBM.SYSDUMMY1
           END-EXEC.

           IF SQLCODE = 0
              CONTINUE
           ELSE
              MOVE SQLCODE                  TO W100-APPO-SQLCODE
              PERFORM   IMPO-ERRO-X-KEY-RANDOM
                 THRU F-IMPO-ERRO-X-KEY-RANDOM
           END-IF.

R14217 F-TROVA-KEY-RANDOM.
           EXIT.
      *================================================================*
       IMPO-ERRO-X-KEY-RANDOM.
      *
           MOVE SPACES                      TO YPCWS001-RIGA
           STRING 'CALCOLO KEY RANDOM ERRATO - '
                  'SQLCODE = ' W100-APPO-SQLCODE
                  DELIMITED BY SIZE
                  INTO YPCWS001-RIGA
           END-STRING
           PERFORM SCRIVI-ST              THRU EX-SCRIVI-ST
      *
           MOVE SPACES                      TO YPCWS001-RIGA
           STRING 'PGPF-PAYEMT-UID ==> '   PGPF-PAYEMT-UID
                  DELIMITED BY SIZE
                  INTO YPCWS001-RIGA
           END-STRING
           PERFORM SCRIVI-ST              THRU EX-SCRIVI-ST
      *
           MOVE SPACES                      TO YPCWS001-RIGA
           PERFORM SCRIVI-ST              THRU EX-SCRIVI-ST
           .
      *
       F-IMPO-ERRO-X-KEY-RANDOM.
           EXIT.
      *================================================================*
R15420 CHIAMA-Z3BCUI99.
      *
           INITIALIZE                    Z3CLUI99
DBG==>*    DISPLAY 'CHIAMA-Z3BCUI99 '
      *
           MOVE 'AUTO'                   TO Z3CLUI99-CANALE
           MOVE 'INQ'                    TO Z3CLUI99-TIPO-RICHIESTA
      *
           MOVE  '2'                     TO Z3CLUI99-TIPO-ID-DISP
           MOVE Z3CLUIFA-OU-PAN-II-TR    TO Z3CLUI99-ID-DISP
DBG==>*    DISPLAY 'Z3CLUI99-ID-DISP       ('Z3CLUI99-ID-DISP')'
           CALL  Z3BCUI99        USING Z3CLUI99
DBG==>*    DISPLAY 'Z3CLUI99-CODI-ERR      ('Z3CLUI99-CODI-ERR')'
DBG==>*    DISPLAY 'Z3CLUI99-FLAG-TIPO-BLOC('Z3CLUI99-FLAG-TIPO-BLOC')'
      *---  Errore generico
           IF Z3CLUI99-CODI-ERR  NOT = '000'
              MOVE SPACES                 TO YPCWS001-RIGA
              STRING 'ERR. ROUTINE Z3BCUI99 '
               'CODI-ERR ' Z3CLUI99-CODI-ERR
                DELIMITED BY SIZE
                INTO YPCWS001-RIGA
              END-STRING
              PERFORM SCRIVI-ST         THRU EX-SCRIVI-ST
              SET WS-SCRI-SCAR-SI         TO TRUE
              SET WS-SALT-ELAB-SI         TO TRUE
      *-
      *--* Imposta area x messaggio errori via mail
              MOVE 'Errore generico routine Z3BCUI99   -'
                                            TO WS-AREA-APPO-YPOE-DESC
              GO TO F-CHIAMA-Z3BCUI99
           END-IF.
           IF Z3CLUI99-FLAG-TIPO-BLOC NOT = 'D'
               MOVE SPACES                  TO YPCWS001-RIGA
               STRING ' - Carta in blocco C0 - '
                         ' Ret.code Uifa('
                       Z3CLUIFA-OU-RET-CODE
                          ')'
                       ' Tipo blocco Ui99('
                       Z3CLUI99-FLAG-TIPO-BLOC
                          ')'
                     DELIMITED BY SIZE
                     INTO YPCWS001-RIGA
               END-STRING
               PERFORM SCRIVI-ST          THRU EX-SCRIVI-ST
               SET WS-SCRI-SCAR-SI          TO TRUE
               SET WS-SALT-ELAB-SI          TO TRUE
      *-
      *--* Imposta area x messaggio errori via mail
               MOVE 'Carta non attiva                   -'
                                            TO WS-AREA-APPO-YPOE-DESC
            END-IF
            .
           IF Z3CLUI99-FLAG-TIPO-BLOC = 'D'
           AND PGPF-DB-CR-FLAG = 'D'
               MOVE SPACES                  TO YPCWS001-RIGA
               STRING ' - Carta in blocco DARE '
                         ' Ret.code Uifa('
                       Z3CLUIFA-OU-RET-CODE
                          ')'
                       ' Tipo blocco Ui99('
                       Z3CLUI99-FLAG-TIPO-BLOC
                          ')'
                     DELIMITED BY SIZE
                     INTO YPCWS001-RIGA
               END-STRING
               PERFORM SCRIVI-ST          THRU EX-SCRIVI-ST
               SET WS-SCRI-SCAR-SI          TO TRUE
               SET WS-SALT-ELAB-SI          TO TRUE
      *-
      *--* Imposta area x messaggio errori via mail
               MOVE 'Carta non attiva                   -'
                                            TO WS-AREA-APPO-YPOE-DESC
            END-IF
            .
R15420 F-CHIAMA-Z3BCUI99.
            EXIT.
      *================================================================*
R11422 TRATTA-CC-BANCARIO.
      *
           IF WS-LETTO-FAD-NO
             PERFORM SELE-TABE-GEP-FAD
                THRU F-SELE-TABE-GEP-FAD
             SET  WS-LETTO-FAD-SI TO TRUE
           END-IF
      *
           PERFORM IMPO-DCD-DEF-DARE    THRU F-IMPO-DCD-DEF-DARE
           PERFORM IMPO-DCD-DARE        THRU F-IMPO-DCD-DARE
      *
      *
           PERFORM INSE-YPTBFAS2
              THRU F-INSE-YPTBFAS2
           IF SQLCODE = ZERO
               PERFORM WRITE-OUTDCD
                  THRU F-WRITE-OUTDCD
               PERFORM IMPOSTA-DATI-FXM2
                  THRU EX-IMPOSTA-DATI-FXM2
               PERFORM SCRIVI-FXM2
                  THRU EX-SCRIVI-FXM2
              SET WS-ELAB-CC-BANCARI-SI        TO TRUE
           ELSE
               PERFORM SCRIVI-SCARTI
                  THRU EX-SCRIVI-SCARTI
           END-IF
            .
R11422 EX-TRATTA-CC-BANCARIO.
            EXIT.
      *================================================================*
R11422 IMPOSTA-DATI-T-FXM2.
      *
           INITIALIZE YPCRREQX
           INITIALIZE OUTFXML-REC
      *
           MOVE WS-NUM-OPER-CC-B      TO YPCRREQX-C-NUM-TRAN
           MOVE WS-TOTALI-IMPO-CC-B   TO YPCRREQX-C-TOTALE-IMP
           MOVE  '0'                  TO YPCRREQX-C-TIPO-REC
           MOVE  SPACES               TO YPCRREQX-DESC
           MOVE  YPCRREQX             TO OUTFXML-REC.
      *
R11422 EX-IMPOSTA-DATI-T-FXM2.
           EXIT.
      *================================================================*
R11422 IMPOSTA-DATI-FXM2.
      *
           COMPUTE WS-TOTALI-IMPO-CC-B =
                   WS-TOTALI-IMPO-CC-B + (PGPF-PAYMT-TOT / 100)
           COMPUTE WS-NUM-OPER-CC-B = WS-NUM-OPER-CC-B + 1
           INITIALIZE YPCRREQX
           INITIALIZE OUTFXM2-REC
           MOVE  YPCWFAMI-O-IBAN-MANDATO TO YPCRREQX-D-IBAN-DEST
           COMPUTE YPCRREQX-D-IMPO-MOV = PGPF-PAYMT-TOT / 100
           MOVE PGPF-PAYEMT-UID          TO   YPCRREQX-D-PAYMENT-UID
           MOVE YPCWFAMI-O-ID-MANDATO    TO   YPCRREQX-D-ID-MANDATO
           MOVE YPCWFAMI-O-DATA-ATTIV-MANDATO TO
                                         YPCRREQX-D-DATA-ATT-MANDATO
      *
R13519*--* Nel caso il rapporto non esista, in base dati del FA,
R13519*--* nel campo ragione sociale viene impostata il valore
R13519*--* del campo insegna, passato da SIA nel flusso PGPF
           IF YPCWFAMI-O-NUME-RAPP-FA = SPACES OR LOW-VALUE
              IF PGPF-ACCT-OWNER-NAM  > SPACE
                 MOVE PGPF-ACCT-OWNER-NAM  TO YPCRREQX-D-RAGI-SOC
              ELSE
                 MOVE PGPF-NAME            TO YPCRREQX-D-RAGI-SOC
              END-IF
           ELSE
              PERFORM CHIAMA-ANAGRAFE-D
              THRU EX-CHIAMA-ANAGRAFE-D
R13519     END-IF
      *
           IF PGPF-BRAND-CODE > SPACE
              MOVE PGPF-BRAND-CODE            TO W100-BRAND-CODE
                                                 W100-BRAND-CODE-CH
              MOVE ' Brand '                  TO W100-DESC-BRAND
                                                 W100-DESC-BRAND-CH
           ELSE
              MOVE SPACE                      TO W100-BRAND-CODE
                                                 W100-BRAND-CODE-CH
              MOVE SPACE                      TO W100-DESC-BRAND
                                                 W100-DESC-BRAND-CH
           END-IF
           IF PGPF-NICKNAME > SPACE
              MOVE PGPF-NICKNAME              TO W100-NICKNAME
                                                 W100-NICKNAME-CH
           ELSE
              MOVE SPACE                      TO W100-NICKNAME
                                                 W100-NICKNAME-CH
           END-IF
           STRING PGPF-BUSINESS-DATE(7:2) '-'
                  PGPF-BUSINESS-DATE(5:2) '-'
                  PGPF-BUSINESS-DATE(1:4)
                  DELIMITED BY SIZE
                  INTO W100-DATA-OPERAZ
           END-STRING
           MOVE PGPF-ME-ID-CODE           TO W100-MERCHANT
                                             W100-MERCHANT-CH
           IF PGPF-FUNCT-CODE   = 200
              MOVE W100-DESCMOV-200        TO WS-DESC-MOV
            ELSE
              MOVE W100-DESCMOV            TO WS-DESC-MOV
           END-IF
           MOVE  WS-DESC-MOV               TO YPCRREQX-DESC
      *
           MOVE  '1'                  TO YPCRREQX-D-TIPO-REC
           MOVE  YPCRREQX             TO OUTFXML-REC.
      *
R11422 EX-IMPOSTA-DATI-FXM2.
           EXIT.
      *==============================================================*
R11422 IMPO-DCD-DEF-DARE.
      *
           INITIALIZE VG0000R
           MOVE 01                        TO VG000-COD-SOC
           MOVE '00'                      TO VG000-COD-UFF
           MOVE '00'                      TO VG000-COD-UFF-ORIG
           MOVE TF03-DESCRIZIONE          TO VG000-DES-TT-MOV-PAR
           PERFORM R315-INCREMENTA-PROGR  THRU R315-INCREMENTA-PROGR-EX
           MOVE W100-PROGR-APERTURA       TO VG000-CNT-PRG-MOV-PAR
           PERFORM R320-CONVERTI-ORA      THRU R320-CONVERTI-ORA-EX
           MOVE W100-DATA-SOLARE-AAMMGG   TO VG000-DAT-CONTABILE-AUTO
      *
           MOVE ZEROES                    TO VG000-DAT-SCADENZA
           MOVE ZEROES                    TO VG000-DAT-CHD
           MOVE '00'                      TO VG000-COD-UFF-NEW-CHS
           MOVE ZEROES                    TO VG000-DAT-VAL
           MOVE ZEROES                    TO VG000-DAT-CAR
           MOVE ZEROES                    TO VG000-DAT-EMI-EFF
           MOVE ZEROES                    TO VG000-DAT-SCD-EFF
           MOVE 'EUR'                     TO VG000-COD-DIVISA
           COMPUTE VG000-IMP-MOV = PGPF-PAYMT-TOT
           MOVE ZEROES                    TO VG000-IMP-MOV-C
           MOVE SPACES                    TO VG000-COD-PRV-SLV
           MOVE YPCWFAMI-O-NUME-RAPP-FA   TO VG000-KEY-PROC(1:12)
           MOVE YPCWFAMI-O-ID-MANDATO     TO VG000-KEY-PROC(13:35)
           .
R11422 F-IMPO-DCD-DEF-DARE.
           EXIT.
      *==============================================================*
R11422 IMPO-DCD-DARE.
           MOVE 'D'                       TO VG000-FLG-SGN
           MOVE W100-DATA-SOLARE9         TO VG000-DAT-SOL
           MOVE TF03-DESCRIZIONE        TO VG000-DES-TT-MOV-PAR
           MOVE TF03-D-DCD-ENTE-4LIV    TO VG000-COD-ENTE-4LIV
           MOVE TF03-D-DCD-TIP-PART     TO VG000-COD-TIP-PART
           MOVE TF03-D-DCD-ENTE-4LIV    TO VG000-COD-ENTE-4LIV-ORIG
           MOVE TF03-D-DCD-ENTE-4LIV    TO VG000-COD-4LI-NEW-CHS
           MOVE TF03-D-DCD-FLG-SEL-TIP-OPE
                                        TO VG000-FLG-SEL-TIP-OPE
           MOVE TF03-D-DCD-PRV-LAV      TO VG000-COD-PRV-LAV
           MOVE TF03-D-DCD-COD-CONT     TO VG000-COD-CONT
           MOVE W100-DATA-SOLARE        TO VG000-DAT-CONTABILE-X8
           .
R11422 F-IMPO-DCD-DARE.
           EXIT.
      *-----------------------------------------------------------------
R11422 SELE-TABE-GEP-FAD.
      *
           INITIALIZE                    HV-TABE
           MOVE 'GEP'                 TO HV-TABE-KNAMTB1
           MOVE 'FAD03'               TO HV-TABE-KVARTB1

           EXEC SQL SELECT DATI
                INTO :HV-TABE-DATI
                FROM XYTBTABE
           WHERE KNAMTB1 = :HV-TABE-KNAMTB1
             AND KVARTB1 = :HV-TABE-KVARTB1
           END-EXEC
           .

           MOVE SQLCODE                  TO W100-APPO-SQLCODE
           EVALUATE SQLCODE
            WHEN +0
                   MOVE HV-TABE-DATI-A  TO YPCRTF03-DATI
           WHEN OTHER
              MOVE SPACES                      TO YPCWS001-RIGA
              MOVE SPACES                      TO YP-MSGERR
              STRING
                  'ERRORE read TABELLA GEP-FAD03 SQLCODE: '
                  W100-APPO-SQLCODE ' riga =>'
              DELIMITED BY SIZE              INTO YP-MSGERR-1
              END-STRING
              MOVE HV-TABE                     TO YP-MSGERR-2
      *
              PERFORM GEST-ERRO-SU-TRE-RIGH
                 THRU F-GEST-ERRO-SU-TRE-RIGH
           END-EVALUATE.
      *
R11422 F-SELE-TABE-GEP-FAD.
           EXIT.
      *-----------------------------------------------------------------
R11422 ROUT-DATE.
      *
           ACCEPT WS-DATA  FROM DATE.
           MOVE WS-DATA                   TO W100-DATA-SOLARE-AAMMGG.
      *
           MOVE W100-DATA-SOLARE          TO W100-DATA-SOLARE9
            .
      *
R11422 F-ROUT-DATE.
           EXIT.
      *-----------------------------------------------------------------
      *================================================================*
R11422 INSE-YPTBFAS2.
      *
           INITIALIZE                    YPDCFAS2
           MOVE YPCWFAMI-O-NUME-RAPP-FA   TO
                                          YPDCFAS2-NUME-RAPP-FA
           EVALUATE YPCWFAMI-O-IBAN-MAND-TIPO-SP
              WHEN  '01'
                     MOVE 'MM'            TO YPDCFAS2-PARTITARIO
              WHEN  '02'
                     MOVE 'CC'            TO YPDCFAS2-PARTITARIO
              WHEN  '03'
                     MOVE 'BA'            TO YPDCFAS2-PARTITARIO
           END-EVALUATE
      *
           MOVE PGPF-PAYMT-TOT            TO  YPDCFAS2-IMPORTO
      *
           MOVE PGPF-PAYEMT-UID           TO YPDCFAS2-PGPF-PAYMENT-UID
      *
           MOVE PGPF-BANK-ACC-TYP         TO YPDCFAS2-PGPF-BANK-ACC-TYP
      *
           STRING    COM-DATE-TIME-H(1:4) '-'
                     COM-DATE-TIME-H(5:2) '-'
                     COM-DATE-TIME-H(7:2)
                     DELIMITED BY SIZE
                     INTO  YPDCFAS2-PGPF-DATA
           END-STRING.
      *
           MOVE 'DCD'                     TO YPDCFAS2-TIPO-PART
           STRING '20'
                  VG000-DAT-CONTABILE-AUTO(1:2) '-'
                  VG000-DAT-CONTABILE-AUTO(3:2) '-'
                  VG000-DAT-CONTABILE-AUTO(5:2)   DELIMITED BY SIZE
                                        INTO YPDCFAS2-DCD-DATA-CONT-AP
           END-STRING
      *
           MOVE '0001-01-01'             TO YPDCFAS2-DCD-DATA-CONT-CH
           MOVE '0001-01-01'             TO YPDCFAS2-DCD-DATA-CONT-AP-KO
           MOVE '0001-01-01'             TO YPDCFAS2-DCD-DATA-CONT-CH-KO
           MOVE '0001-01-01-00.00.00.000000'
                                         TO YPDCFAS2-TMST-DISTINTA
           MOVE '0001-01-01-00.00.00.000000'
                                         TO YPDCFAS2-TMST-CHIUSURA
           MOVE '0001-01-01-00.00.00.000000'
                                         TO YPDCFAS2-TMST-ESITO
      *
           MOVE VG000-CNT-PRG-MOV-PAR     TO YPDCFAS2-DCD-PROG-PART
           MOVE VG000-COD-TIP-PART        TO YPDCFAS2-DCD-TIPO-PART
           MOVE VG000-FLG-SGN             TO YPDCFAS2-DCD-SEGNO-AP
           MOVE YPCWFAMI-O-ID-MANDATO       TO YPDCFAS2-ID-MANDATO
           IF   YPCWFAMI-O-DATA-ATTIV-MANDATO NOT = SPACE AND
                LOW-VALUE
           MOVE YPCWFAMI-O-DATA-ATTIV-MANDATO
                                    TO YPDCFAS2-DATA-ATTIV-MANDATO
           ELSE
                MOVE '0001-01-01'
                                    TO YPDCFAS2-DATA-ATTIV-MANDATO
           END-IF
           MOVE YPCWFAMI-O-IBAN-MANDATO     TO YPDCFAS2-IBAN

            EXEC SQL INSERT INTO YPTBFAS2
            (
             NUME_RAPP_FA
            ,PARTITARIO
            ,IMPORTO
            ,PGPF_PAYMENT_UID
            ,PGPF_BANK_ACC_TYP
            ,PGPF_DATA
            ,TIPO_PART
            ,DCD_DATA_CONT_AP
            ,DCD_DATA_CONT_CH
            ,DCD_DATA_CONT_AP_KO
            ,DCD_DATA_CONT_CH_KO
            ,DCD_PROG_PART
            ,DCD_TIPO_PART
            ,DCD_SEGNO_AP
            ,ID_MANDATO
            ,DATA_ATTIV_MANDATO
            ,IBAN
            ,TMST_INSE
            ,TMST_DISTINTA
            ,TMST_CHIUSURA
            ,TMST_ESITO
             )
            VALUES (
              :YPDCFAS2-NUME-RAPP-FA
             ,:YPDCFAS2-PARTITARIO
             ,:YPDCFAS2-IMPORTO
             ,:YPDCFAS2-PGPF-PAYMENT-UID
             ,:YPDCFAS2-PGPF-BANK-ACC-TYP
             ,:YPDCFAS2-PGPF-DATA
             ,:YPDCFAS2-TIPO-PART
             ,:YPDCFAS2-DCD-DATA-CONT-AP
             ,:YPDCFAS2-DCD-DATA-CONT-CH
             ,:YPDCFAS2-DCD-DATA-CONT-AP-KO
             ,:YPDCFAS2-DCD-DATA-CONT-CH-KO
             ,:YPDCFAS2-DCD-PROG-PART
             ,:YPDCFAS2-DCD-TIPO-PART
             ,:YPDCFAS2-DCD-SEGNO-AP
             ,:YPDCFAS2-ID-MANDATO
             ,:YPDCFAS2-DATA-ATTIV-MANDATO
             ,:YPDCFAS2-IBAN
             , CURRENT TIMESTAMP
             ,:YPDCFAS2-TMST-DISTINTA
             ,:YPDCFAS2-TMST-CHIUSURA
             ,:YPDCFAS2-TMST-ESITO
                                )
           END-EXEC.
      *
           MOVE SQLCODE                  TO W100-APPO-SQLCODE
           IF SQLCODE = 0
              ADD 1                         TO CTR-TABFAS2-INSE
           ELSE
              IF SQLCODE = -803
                 ADD 1                      TO W300-SCART-DUPKEY
                 PERFORM   IMPO-ERRO-X-DUP-KEY-FAS2
                    THRU F-IMPO-ERRO-X-DUP-KEY-FAS2
              ELSE
                 PERFORM   IMPO-ERRO-INSE-FAS2
                    THRU F-IMPO-ERRO-INSE-FAS2
              END-IF
           END-IF.
      *
R11422 F-INSE-YPTBFAS2.
           EXIT.
      *================================================================*
      * WRITE FLUSSO DC2                                               *
      *================================================================*
R11422 STMP-RIGH-T11.
      *
           MOVE    SPACES     TO        YPCWS001-RIGA
           STRING  'TOTALE RECORDS FILE FXM2:           '
                   ETR-CONT-FXM2
                   DELIMITED BY SIZE
                   INTO YPCWS001-RIGA
           END-STRING
           PERFORM SCRIVI-ST  THRU      EX-SCRIVI-ST
           .
R11422 F-STMP-RIGH-T11.
           EXIT.
      *================================================================*
R11422 STMP-RIGH-T12.
      *
           MOVE    SPACES     TO        YPCWS001-RIGA
           PERFORM SCRIVI-ST  THRU      EX-SCRIVI-ST
      *
           MOVE    SPACES     TO        YPCWS001-RIGA
           STRING  'RIGHE INSERITE IN TABELLA YPTBFAS2_:'
                  ETR-TABFAS2-INSE
                  DELIMITED BY SIZE
                  INTO  YPCWS001-RIGA
           END-STRING.
           PERFORM SCRIVI-ST THRU EX-SCRIVI-ST
           .
R11422 F-STMP-RIGH-T12.
           EXIT.
      *================================================================*
R11422 CHIAMA-ANAGRAFE-D.
      *
           INITIALIZE    AREA-ACS108A
           INITIALIZE                        L-ACS108-ARG
           MOVE ZEROES                    TO L-ACS108-I-BANCA
           MOVE ' '                       TO L-ACS108-I-TIPO-RICH
           MOVE ZEROES                    TO L-ACS108-I-DATA-RIF
           MOVE 'FA '                     TO L-ACS108-I-SERVIZIO
           MOVE ZEROES                    TO WK-END                     00077200
           MOVE YPCWFAMI-O-NUME-RAPP-FA  TO L-ACS108-I-NUMERO
           MOVE YPCWFAMI-O-FILI-RAPP     TO L-ACS108-I-FILIALE
           CALL WK-ACS108BT             USING L-ACS108-ARG.
           .
           IF  L-ACS108-RET-CODE = ZEROES
               IF    L-ACS108-COGNOME  = SPACES AND
                     L-ACS108-NOME  = SPACES
                    MOVE  L-ACS108-RAGSOC-1      TO YPCRREQX-D-RAGI-SOC
               ELSE
                   MOVE  L-ACS108-COGNOME        TO
                                        YPCRREQX-D-RAGI-SOC(1:30)
                   MOVE  L-ACS108-NOME           TO
                                        YPCRREQX-D-RAGI-SOC(31:)
               END-IF
               IF  L-ACS108-PARTITA-IVA  NOT = SPACE
                    MOVE  L-ACS108-PARTITA-IVA TO
                                          YPCRREQX-D-CODICE-FISC-DEB
               ELSE
????               MOVE  L-ACS108-COD-FISCALE   TO
                                          YPCRREQX-D-CODICE-FISC-DEB
               END-IF
               MOVE  L-ACS108-IND-SEDE-LEG   TO YPCRREQX-D-INDIRIZZO
               MOVE  L-ACS108-CAP-SEDE-LEG   TO YPCRREQX-D-CAP
               MOVE  L-ACS108-LOC-SEDE-LEG   TO YPCRREQX-D-LOC
               MOVE  L-ACS108-PROV-SEDE-LEG  TO YPCRREQX-D-PROV
               MOVE  L-ACS108-NAZ-SEDE-LEG   TO YPCRREQX-D-NAZ
            .
      *
R11422 EX-CHIAMA-ANAGRAFE-D.
           EXIT.
      *================================================================*
TK1274 CONTROLLA-GEP-FPR.
      *
           SET TROVATO-SU-FPR-NO  TO TRUE
           INITIALIZE                    HV-TABE
           MOVE 'GEP'                 TO HV-TABE-KNAMTB1
           MOVE 'FPR'                 TO HV-TABE-KVARTB1(1:3)
           MOVE YPCWFAMI-O-PROD-E-PROD(1)
                                      TO HV-TABE-KVARTB1(4:12)

DBG==>*    display 'CONTROLLA-GEP-FPR'
DBG==>*    display 'HV-TABE-KVARTB1 : ' HV-TABE-KVARTB1
           EXEC SQL SELECT DATI
                INTO :HV-TABE-DATI
                FROM XYTBTABE
           WHERE KNAMTB1 = :HV-TABE-KNAMTB1 AND
                 KVARTB1 = :HV-TABE-KVARTB1
           END-EXEC
           .

      *
           MOVE SQLCODE                  TO W100-APPO-SQLCODE
DBG==>*    display 'W100-APPO-SQLCODE :' W100-APPO-SQLCODE
           EVALUATE SQLCODE
            WHEN +0
                   MOVE HV-TABE-DATI-A      TO YPCRTFPR-DATI
                   SET TROVATO-SU-FPR-SI  TO TRUE
                   ADD 1 TO CTR-TABFPR-LETTE
DBG==>*    display 'YPCRTFPR-DATI :'  YPCRTFPR-DATI
            WHEN +100
                   SET TROVATO-SU-FPR-NO  TO TRUE
                   ADD 1 TO CTR-TABFPR-NOT-FOUND
            WHEN OTHER
              MOVE SPACES                      TO YPCWS001-RIGA
              MOVE SPACES                      TO YP-MSGERR
              STRING
                  'ERRORE read TABELLA xytbtabe SQLCODE: '
                  W100-APPO-SQLCODE ' riga =>'
              DELIMITED BY SIZE              INTO YP-MSGERR-1
              END-STRING
              MOVE HV-TABE                     TO YP-MSGERR-2
      *
              PERFORM GEST-ERRO-SU-TRE-RIGH
                 THRU F-GEST-ERRO-SU-TRE-RIGH
           END-EVALUATE.
      *
           .
TK1274 EX-CONTROLLA-GEP-FPR.
           EXIT.
